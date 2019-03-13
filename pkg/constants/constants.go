package constants

// Queries
const (
	// Environment variables to find DB username and password
	UserEnvVar = "DB_USER"
	PassEnvVar = "DB_PASS"

	// 20 chars can be got by hex encoding 10 bytes
	SerialLenHexBytes = 20 / 2

	// Create queries
	CreateQuery = `INSERT INTO ims_item
(serial, fulfillment_center_sk, mgmnt_lkp_id, tags, fulfillment_mode, audit_wip, damaged, created_at, updated_at)
VALUES ($1, $2, $3, '{}', 'BS', false, false, NOW(), NOW())`
	BatchCreateQuery = `INSERT INTO ims_item
(serial, fulfillment_center_sk, mgmnt_lkp_id, tags, fulfillment_mode, audit_wip, damaged, created_at, updated_at)
VALUES %s`

	// Update queries
	RTDQuery = `UPDATE ims_item
SET status='RTD', updated_at=NOW()
WHERE serial IN ('%s')`
	RCVQuery = `UPDATE ims_item
SET status='RCV', received_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	HLDQuery = `UPDATE ims_item
SET status='HLD', hold_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	TBQQuery = `UPDATE ims_item
SET status='TBQ', qc_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	ACCQuery = `UPDATE ims_item
SET status='ACC', updated_at=NOW()
WHERE serial IN ('%s')`
	REJQuery = `UPDATE ims_item
SET status='REJ', updated_at=NOW()
WHERE serial IN ('%s')`
	LSTQuery = `UPDATE ims_item
SET status='LST', lost_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	STKQuery = `UPDATE ims_item
SET status='STK', stocked_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	ALCQuery = `UPDATE ims_item
SET status='ALC', allocated_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	TBPQuery = `UPDATE ims_item
SET status='TBP', tbp_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	PIKQuery = `UPDATE ims_item
SET status='PIK', picked_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	PAKQuery = `UPDATE ims_item
SET status='PAK', packed_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	RTSQuery = `UPDATE ims_item
SET status='RTS', rts_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	SHPQuery = `UPDATE ims_item
SET status='SHP', shipped_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	CANQuery = `UPDATE ims_item
SET status='CAN', updated_at=NOW()
WHERE serial IN ('%s')`
	RPOQuery = `UPDATE ims_item
SET status='RPO', rtv_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
	RWRQuery = `UPDATE ims_item
SET status='RWR', returned_at=NOW(), updated_at=NOW()
WHERE serial IN ('%s')`
)

var (
	// QueryMap is a map of status and the corresponding update query
	QueryMap = map[string]string{
		"RTD": RTDQuery,
		"RCV": RCVQuery,
		"HLD": HLDQuery,
		"TBQ": TBQQuery,
		"ACC": ACCQuery,
		"REJ": REJQuery,
		"LST": LSTQuery,
		"STK": STKQuery,
		"ALC": ALCQuery,
		"TBP": TBPQuery,
		"PIK": PIKQuery,
		"PAK": PAKQuery,
		"RTS": RTSQuery,
		"SHP": SHPQuery,
		"CAN": CANQuery,
		"RPO": RPOQuery,
		"RWR": RWRQuery,
	}
)
