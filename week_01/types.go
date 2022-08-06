package main

type UserID int // not int, UserID type, incompatible

func main() {
	idx := 1
	var uid UserID = 42

	//uid = idx // denied
	uid = UserID(idx) // simple cast
	println("idx, uid: ", idx, uid)

}
