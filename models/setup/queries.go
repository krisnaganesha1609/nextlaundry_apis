package setup

func UsersTriggers() {
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_users_insert_history` AFTER INSERT ON `tb_user` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('User with id ', NEW.id, ' (', NEW.nama, ')', ' has been added'), 'users', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_users_update_history` AFTER UPDATE ON `tb_user` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('User with id ', OLD.id, ' (', OLD.nama, ')', ' has been updated'), 'users', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_users_delete_history` AFTER DELETE ON `tb_user` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('User with id ', OLD.id, ' (', OLD.nama, ')', ' has been removed'), 'users', NOW(), NULL);")
}

func OutletTriggers() {
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_outlet_insert_history` AFTER INSERT ON `tb_outlet` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Outlet with id ', NEW.id, ' (', NEW.nama, ')', ' has been added'), 'outlet', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_outlet_update_history` AFTER UPDATE ON `tb_outlet` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Outlet with id ', OLD.id, ' (', OLD.nama, ')', ' has been updated'), 'outlet', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_outlet_delete_history` AFTER DELETE ON `tb_outlet` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Outlet with id ', OLD.id, ' (', OLD.nama, ')', ' has been removed'), 'outlet', NOW(), NULL);")
}

func MemberTriggers() {
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_member_insert_history` AFTER INSERT ON `tb_member` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Member with id ', NEW.id, ' (', NEW.nama, ')', ' has been added'), 'member', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_member_update_history` AFTER UPDATE ON `tb_member` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Member with id ', OLD.id, ' (', OLD.nama, ')', ' has been updated'), 'member', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_member_delete_history` AFTER DELETE ON `tb_member` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Member with id ', OLD.id, ' (', OLD.nama, ')', ' has been removed'), 'member', NOW(), NULL);")
}

func PackageTriggers() {
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_paket_insert_history` AFTER INSERT ON `tb_paket` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Package with id ', NEW.id, ' (', NEW.nama_paket, ')', ' has been added'), 'paket', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_paket_update_history` AFTER UPDATE ON `tb_paket` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Package with id ', OLD.id, ' (', OLD.nama_paket, ')', ' has been updated'), 'paket', NOW(), NULL);")
	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_paket_delete_history` AFTER DELETE ON `tb_paket` FOR EACH ROW INSERT INTO log_history VALUES(NULL, CONCAT('Package with id ', OLD.id, ' (', OLD.Nama_paket, ')', ' has been removed'), 'paket', NOW(), NULL);")
}

// func TransactionTriggers() {
// 	DB.Exec("CREATE OR REPLACE TRIGGER `update` AFTER INSERT ON `tb_transaksi` FOR EACH ROW CALL InsertFirstTransactionLog();")
// 	DB.Exec("CREATE OR REPLACE TRIGGER `trigger_detailtransaction_insert_history` AFTER INSERT ON `tb_detail_transaksi` FOR EACH ROW CALL UpdateTheTransactionLog();")
// }
