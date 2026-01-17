package repository

const (
	userCreate      = `INSERT INTO "User" (username,email,password_hash,user_type) VALUES ($1,$2,$3,$4) RETURNING *`
	userFind        = `SELECT * FROM "User" WHERE id = $1 LIMIT 1`
	userFindByEmail = `SELECT * FROM "User" WHERE email = $1 LIMIT 1`
	userUpdate      = `UPDATE "User" SET
username = COALESCE(NULLIF($2, ''), username),
email = COALESCE(NULLIF($3, ''), email),
updated_at = NOW() WHERE id = $1 RETURNING *`
	userDelete = `DELETE FROM "User" WHERE id = $1`

	userFindByEmailWithSeller = `
		SELECT u.*, s.id as seller_id, s.seller_name
		FROM "User" u
		INNER JOIN "Seller" s ON u.id = s.user_id
		WHERE u.email = $1
		LIMIT 1`

	userFindByEmailWithCollector = `
		SELECT u.*, c.id as collector_id, c.collector_name
		FROM "User" u
		INNER JOIN "Collector" c ON u.id = c.user_id
		WHERE u.email = $1
		LIMIT 1`

	userFindByEmailWithCompany = `
		SELECT u.*, co.id as company_id, co.company_name
		FROM "User" u
		INNER JOIN "Company" co ON u.id = co.user_id
		WHERE u.email = $1
		LIMIT 1`

	sellTransactionCreate = `INSERT INTO "SellTransaction" (seller_id, collector_id, volume, price)
		VALUES ($1, $2, $3, $4) RETURNING *`

	distributeTransactionCreate = `INSERT INTO "DistributeTransaction" (collector_id, company_id, volume, price)
		VALUES ($1, $2, $3, $4) RETURNING *`

	sellTransactionUpdate = `UPDATE "SellTransaction" SET volume = $2, price = $3, updated_at = NOW()
		WHERE id = $1 RETURNING *`

	distributeTransactionUpdate = `UPDATE "DistributeTransaction" SET volume = $2, price = $3, updated_at = NOW()
		WHERE id = $1 RETURNING *`

	sellTransactionFindById = `SELECT * FROM "SellTransaction" WHERE id = $1 LIMIT 1`

	distributeTransactionFindById = `SELECT * FROM "DistributeTransaction" WHERE id = $1 LIMIT 1`

	updateCollectorVolume = `UPDATE "Collector" SET total_volume = COALESCE(total_volume, 0) + $2, updated_at = NOW()
		WHERE id = $1`

	reportSalesByDate = `SELECT
		st.created_at as transaction_date,
		c.collector_name,
		s.seller_name,
		st.volume as oil_volume,
		st.price
	FROM "SellTransaction" st
	JOIN "Seller" s ON st.seller_id = s.id
	JOIN "Collector" c ON st.collector_id = c.id
	WHERE st.created_at >= $1 AND st.created_at <= $2
	ORDER BY st.created_at DESC`

	reportPurchasesByDate = `SELECT
		dt.created_at as transaction_date,
		c.collector_name,
		co.company_name,
		dt.volume as oil_volume,
		dt.price
	FROM "DistributeTransaction" dt
	JOIN "Collector" c ON dt.collector_id = c.id
	JOIN "Company" co ON dt.company_id = co.id
	WHERE dt.created_at >= $1 AND dt.created_at <= $2
	ORDER BY dt.created_at DESC`

	reportAllSales = `SELECT
		st.created_at as transaction_date,
		c.collector_name,
		s.seller_name,
		st.volume as oil_volume,
		st.price
	FROM "SellTransaction" st
	JOIN "Seller" s ON st.seller_id = s.id
	JOIN "Collector" c ON st.collector_id = c.id
	ORDER BY st.created_at DESC`

	reportAllPurchases = `SELECT
		dt.created_at as transaction_date,
		c.collector_name,
		co.company_name,
		dt.volume as oil_volume,
		dt.price
	FROM "DistributeTransaction" dt
	JOIN "Collector" c ON dt.collector_id = c.id
	JOIN "Company" co ON dt.company_id = co.id
	ORDER BY dt.created_at DESC`

	oilCreate = `INSERT INTO "Oil" (collector_id, total_volume)
		VALUES ($1, 0) RETURNING *`

	oilFind = `SELECT * FROM "Oil" WHERE id = $1 LIMIT 1`

	oilGetByCollectorId = `SELECT * FROM "Oil" WHERE collector_id = $1 LIMIT 1`

	oilUpdate = `UPDATE "Oil" SET
		total_volume = $2,
		updated_at = NOW()
		WHERE id = $1 RETURNING *`

	oilDelete = `DELETE FROM "Oil" WHERE id = $1`
)
