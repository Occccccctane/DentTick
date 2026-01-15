package Domain

type User struct {
	Id       int64
	Name     string
	Info     string
	Password string
	//医生 Id和患者 Id互斥,身份由Identity字段决定
	//1为医生，2为患者,3为管理员
	Identity  int8
	DoctorId  string
	PatientId string
	Phone     string
}
