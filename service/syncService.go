package service

import (
	"context"
	"fmt"
	"strings"

	models "github.com/cs161079/godbLib/Models"
	utils "github.com/cs161079/godbLib/Utils"
	logger "github.com/cs161079/godbLib/Utils/goLogger"
	"github.com/cs161079/godbLib/db"
	"github.com/cs161079/godbLib/repository"
	"github.com/cs161079/godbLib/service"
	"gorm.io/gorm"
)

type SyncService interface {
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα για τις γραμμές των λεοφωρείων από τον Server του OASA στην δική
	// μας βάση δεδομένων. Καλούμε το API /webGetLinesWithMLInfo το οποίο μας επιστρέφει όλες τις γραμμές με την μορφή
	// JSON, δηλαδή ένα πίνακα από record. Η κάθε γραμμής περιέχει την εξής πληροφορία
	//	{
	//	    "ml_code": "9",
	//	    "sdc_code": "54",
	//	    "line_code": "1151",
	//	    "line_id": "021",
	//	    "line_descr": "ΠΛΑΤΕΙΑ ΚΑΝΙΓΓΟΣ - ΓΚΥΖH (ΚΥΚΛΙΚΗ)",
	//	    "line_descr_eng": "PLATEIA KANIGKOS - GKIZI",
	//	    "mld_master": "1"
	//	}
	// =================================================================================================================
	SyncLines(context.Context) error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των διαδρομών από τον Server του OASA στην δική μας βάση δεδομένων.
	// Καλούμε το API /getRoutes το οποίο μας επιστρέφει όλες τις διαδρομές σε txt μορφή, τα δεδομένα των διαδρομών
	// χωρισμένα με κόμμα και κάθε διαδρομή από την άλλη χωρίζονται με κόμμα. Αυτά είναι τα δεδομένα μίας διαδρομής.
	//
	// (1754,799, "ΕΛ.ΒΕΝΙΖΕΛΟΥ - ΚΑΙΣΑΡΙΑΝΗ", "EL. VENIZELOU - KAISARIANI",2,9889.61)
	// =================================================================================================================
	SyncRoutes(context.Context) error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των στάσεων από τον Server του OASA στην δική μας βάση δεδομένων.
	// Καλούμε το API /getStops το οποίο μας επιστρέφει όλες τις στάσεις σε txt μορφή, τα δεδομένα των στάσεων
	// χωρισμένα με κόμμα και κάθε στάση από την άλλη χωρίζονται ομοιώς με κόμμα. Αυτά είναι τα δεδομένα μίας στάσης.
	//
	// (10001, "010001", "ΣΤΡΟΦΗ", "STROFH", "ΕΛ.ΒΕΝΙΖΕΛΟΥ", "ΕΛ.ΒΕΝΙΖΕΛΟΥ", -1,23.665,37.9986,0,0,
	//                       "| ΑΝΩ ΑΓ. ΒΑΡΒΑΡΑ| ΠΕΙΡΑΙΑΣ ΠΛ. ΚΑΡΑΪΣΚΑΚΗ", "| ANO AG. BARBARA| PEIRAIAS PL. KARAISKAKΗ")
	// =================================================================================================================
	SyncStops(context.Context) error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των στάσεων ανά διαδρομή από τον Server του OASA στην δική μας
	// βάση δεδομένων. Καλούμε το API /getRouteStops το οποίο μας επιστρέφει όλες τις στάσεις σε txt μορφή, τα δεδομένα
	// των στάσεων ανά διαδρομή είναι χωρισμένα με κόμμα και κάθε γραμμή είναι και μία εγγραφή στην βάση.
	// Αυτά είναι τα δεδομένα μίας εγγραφής.
	//
	//	(103406,2081,10373,1)
	// =================================================================================================================
	SyncRouteStops(context.Context) error
}

type syncService struct {
	LineSrv service.LineService
}

func NewSyncService() SyncService {
	return syncService{}
}

func recPreparation(recStr string) string {
	var trimmedSpace = strings.ReplaceAll(recStr, " ", "")
	return strings.ReplaceAll(trimmedSpace, "\"", "")
}

func (s syncService) SyncLines(ctx context.Context) error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************

	s.LineSrv = service.NewLineService(repository.NewLineRepository(dbConnection))
	var restSrv = service.NewRestService()

	response := restSrv.OasaRequestApi00("webGetLinesWithMLInfo", nil)
	if response.Error != nil {
		return response.Error
	}
	txt := dbConnection.Begin()
	if err := s.LineSrv.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
	}
	logger.INFO("Delete all data from Line table in database succesfully.")
	var lineArray []models.Line = make([]models.Line, 0)
	logger.INFO("Start sychronize data from OASA Server...")
	for _, ln := range response.Data.([]any) {
		lineOasa := s.LineSrv.GetMapper().GeneralLine(ln.(map[string]interface{}))
		line := s.LineSrv.GetMapper().OasaToLine(lineOasa)

		lineArray = append(lineArray, line)
		if len(lineArray) == 1000 {
			_, err := s.LineSrv.WithTrx(txt).InsertArray(lineArray)
			if err != nil {
				txt.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(lineArray)))
			lineArray = make([]models.Line, 0)
		}

	}

	if len(lineArray) > 0 {
		_, err := s.LineSrv.WithTrx(txt).InsertArray(lineArray)
		if err != nil {
			txt.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(lineArray)))
	}

	txt.Commit()
	logger.INFO("Finished sychronize data from OASA Server.")
	return nil
}

func (s syncService) SyncRoutes(ctx context.Context) error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = service.NewRestService()

	var routeSrv = service.NewRouteService(dbConnection)
	response := restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		return response.Error
	}

	tx := dbConnection.Begin()

	err := routeSrv.WithTrx(tx).DeleteAll()
	if err != nil {
		tx.Rollback()
		return err
	}
	var routeArray []models.Route = make([]models.Route, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των διαδρομών **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 6 {
			return fmt.Errorf("Η γραμμή του Record  είναι ελλειπής.")
		}
		rt := models.Route{}
		num, err := utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.Route_Code = *num
		num, err = utils.StrToInt32(recordArr[1])
		if err != nil {
			return err
		}
		rt.Line_Code = *num
		rt.Route_Descr = recordArr[2]
		rt.Route_Descr_eng = recordArr[3]
		num, err = utils.StrToInt32(recordArr[4])
		if err != nil {
			return err
		}
		rt.Route_Type = int8(*num)
		fl32 := utils.StrToFloat32(recordArr[5])
		rt.Route_Distance = fl32

		routeArray = append(routeArray, rt)

		if len(routeArray) == 10000 {
			// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
			_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
		}

	}

	if len(routeArray) > 0 {
		// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
	}
	tx.Commit()

	return nil
}

func (s syncService) SyncStops(ctx context.Context) error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = service.NewRestService()

	var stopSrv = service.NewStopService(repository.NewStopRepository(dbConnection))
	response := restSrv.OasaRequestApi02("getStops")
	if response.Error != nil {
		return response.Error
	}

	tx := dbConnection.Begin()

	err := stopSrv.WithTrx(tx).DeleteAll()
	if err != nil {
		tx.Rollback()
		return err
	}
	logger.INFO("Η διαγραφή των στάσεων έγινε με επιτυχία.")

	var stopArr []models.Stop = make([]models.Stop, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	logger.INFO("Έναρξη συγχρονισμού δεδομένων...")
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των στάσεων **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 13 {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}
		rt := models.Stop{}
		num32, err := utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.Stop_code = *num32
		rt.Stop_id = recordArr[1]
		rt.Stop_descr = recordArr[2]
		rt.Stop_descr_eng = recordArr[3]
		rt.Stop_street = recordArr[4]
		rt.Stop_street_eng = recordArr[5]
		num32, err = utils.StrToInt32(recordArr[6])
		if err != nil {
			return err
		}
		rt.Stop_heading = *num32
		rt.Stop_lng = utils.StrToFloat(recordArr[7])
		rt.Stop_lat = utils.StrToFloat(recordArr[8])
		num8, err := utils.StrToInt8(recordArr[9])
		if err != nil {
			return err
		}
		rt.Stop_type = *num8
		num8, err = utils.StrToInt8(recordArr[10])
		if err != nil {
			return err
		}
		rt.Stop_amea = *num8
		rt.Destinations = recordArr[11]
		rt.Destinations_Eng = recordArr[12]

		stopArr = append(stopArr, rt)

		if len(stopArr) == 1000 {
			// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
			_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
			stopArr = make([]models.Stop, 0)
		}

	}

	if len(stopArr) > 0 {
		// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
	}

	tx.Commit()

	return nil
}

func (s syncService) SyncRouteStops(ctx context.Context) error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = service.NewRestService()

	var routeSrv = service.NewRouteService(dbConnection)
	response := restSrv.OasaRequestApi02("getRouteStops")
	if response.Error != nil {
		return response.Error
	}

	tx := dbConnection.Begin()
	err := routeSrv.WithTrx(tx).DeleteRoute02()
	if err != nil {
		tx.Rollback()
		return err
	}
	logger.INFO("Delete from Route02 Table succesfully.")

	var route02Arr []models.Route02 = make([]models.Route02, 0)
	logger.INFO("Start sychronization Route02 data from OASA Server...")
	for _, rec := range response.Data.([]string) {
		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(4) {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route02{}
		num32, err := utils.StrToInt32(row[1])
		if err != nil {
			return err
		}
		rt.Route_code = *num32
		num64, err := utils.StrToInt64(row[2])
		rt.Stop_code = *num64
		num16, err := utils.StrToInt16(row[3])
		rt.Senu = *num16

		route02Arr = append(route02Arr, rt)
		if len(route02Arr) == 10000 {
			_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
			route02Arr = make([]models.Route02, 0)
		}
	}
	if len(route02Arr) > 0 {
		_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
	}
	tx.Commit()
	logger.INFO("Finished sychronization Route02 data from OASA Server.")
	return nil
}
