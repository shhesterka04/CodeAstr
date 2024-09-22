package db

func (db Database) Close() {
	db.cluster.Close()
}
