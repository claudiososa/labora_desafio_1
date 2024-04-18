package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/manifoldco/promptui"
)

type User struct {
	id       int
	username string
	password string
	name     string
	role     string
}

type UserAuthenticated struct {
	username string
	name     string
	role     string
	status   bool
}

var (
	directorActions = []string{"Listar profesores", "Agregar Profesor", "Eliminar Profesor", "Generar mensaje para profesor"}
	teacherActions  = []string{"Listar alumnos", "Agregar Alumno", "Eliminar Alumno", "Generar mensaje para alumno"}
	laborers        []string
	actionCounter   = make(map[string]int)
)

func main() {

	var userAuthenticated UserAuthenticated

	var users []User

	users = []User{
		{id: 1, username: "juan", password: "juan", name: "Juan Soria", role: "director"},
		{id: 2, username: "carina", password: "carina", name: "Carina Lopez", role: "profesor"},
		{id: 4, username: "andrea", password: "andrea", name: "Andrea Bernal", role: "profesor"},
		{id: 5, username: "pedro", password: "pedro", name: "Pedro Moreno", role: "alumno"},
	}

	for {
		clearScreen()

		switch userAuthenticated.role {

		//Si no hay usuario logueado.
		case "":
			printSimpleText("Seleccione una opcion.")
			printSimpleText("1. Login.")
			printSimpleText("0. Salir.")

			var option int
			fmt.Scanln(&option)

			switch option {
			case 1:
				clearScreen()
				var username string
				var password string

				printSimpleText("Username: ")
				fmt.Scanln(&username)
				printSimpleText("password: ")
				fmt.Scanln(&password)

				userAuthenticated = loginUser(username, password, users)

				if userAuthenticated.status == false {
					clearScreen()
					printSimpleText("ATENCION! Datos de acceso incorrectos")
					printLine()
					printSimpleText("Presione una tecla para continuar...")
					fmt.Scanln()
				}
			case 0:
				printSimpleText("---- Cerrando sesion. Hasta pronto ----")
				os.Exit(0)

			}

		//El usuario logueado es de tipo director
		case "director":
			clearScreen()
			fmt.Println("Nombre del usuario: ", userAuthenticated.name, "/ Role: ", userAuthenticated.role)
			printLine()
			printSimpleText("Seleccione una opcion.")
			printLine()
			for i, action := range directorActions {
				fmt.Printf("%d. %s\n", i+1, action)
			}
			printSimpleText("0. Cerrar sesion")

			var option int
			fmt.Scanln(&option)

			switch option {
			case 1:
				listUsers("profesor", users, false)
			case 2:
				var newTeacher User

				//Obtengo el proximo Id para Users
				nextId := getNextId(users)

				//Obtengo el objeto que contiene la data del nuevo usuario de tipo teacher
				newTeacher = createNewUser(nextId, "profesor")

				//Agrego el nuevo teacher a users
				users = append(users, newTeacher)

				break
			case 3:
				users = deleteUser("profesor", users)
				break
			case 4:
				createNewMessage()
				break
			case 0: //Cerrar Sesion
				userAuthenticated = UserAuthenticated{
					username: "",
					role:     "",
					status:   false,
				}
				clearScreen()
			}

			//El usuario logueado es de tipo profesor
		case "profesor":
			clearScreen()
			fmt.Println("Nombre del usuario: ", userAuthenticated.name, "/ Role: ", userAuthenticated.role)
			printLine()
			printSimpleText("Seleccione una opcion.")
			printLine()
			for i, action := range teacherActions {
				fmt.Printf("%d. %s\n", i+1, action)
			}
			printSimpleText("0. Cerrar sesion")

			var option int
			fmt.Scanln(&option)

			switch option {
			case 1:
				listUsers("alumno", users, false)
			case 2:
				var newStudent User

				//Obtengo el proximo Id para Users
				nextId := getNextId(users)

				//Obtengo el objeto que contiene la data del nuevo usuario de tipo teacher
				newStudent = createNewUser(nextId, "alumno")

				//Agrego el nuevo alumno a users
				users = append(users, newStudent)

				break
			case 3:
				users = deleteUser("alumno", users)
				break

			case 0: //Cerrar Sesion
				userAuthenticated = UserAuthenticated{
					username: "",
					role:     "",
					status:   false,
				}
				clearScreen()
			}
		}
	}

}

func printLine() {
	fmt.Println("/--------------------------------------------------------/")
}

func printSimpleText(text string) {
	fmt.Println(text)
}

func loginUser(username string, password string, users []User) UserAuthenticated {

	//find username
	for _, user := range users {
		if user.username == username && user.password == password {
			return UserAuthenticated{
				username: user.username,
				name:     user.name,
				role:     user.role,
				status:   true,
			}
		}
	}

	return UserAuthenticated{
		username: "",
		name:     "",
		role:     "",
		status:   false,
	}
}

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
	printLine()
	printSimpleText("Sistema Administracioon de Colegio Secundario.")
	printLine()
}

func listUsers(role string, users []User, onlyLists bool) {
	printLine()
	fmt.Println("Listado de ", role)
	printLine()
	for _, user := range users {
		if user.role == role {
			fmt.Println("Name: ", user.name, "Rol: ", user.role, "id", user.id)
		}
	}
	if !onlyLists {
		printLine()
		printSimpleText("Presione una tecla para continuar...")
		fmt.Scanln()
	}

}

func getNextId(users []User) int {
	var maxId int

	//encuentra el ultimo id de los usuarios para agregarle un +1 al nuevo creado y evitar duplicados
	for _, user := range users {
		if user.id >= maxId {
			maxId = user.id + 1
		}
	}
	return maxId
}

func createNewUser(nextId int, userType string) User {
	clearScreen()

	var name string

	printSimpleText("Nombre Completo:")
	fmt.Scanln(&name)

	teacher := User{
		id:       nextId,
		username: name,
		password: name,
		name:     name,
		role:     userType,
	}
	return teacher
}

func deleteUser(userType string, users []User) []User {
	clearScreen()
	listUsers(userType, users, true)

	var idUser int

	printLine()
	fmt.Println("Ingrese el id del ", userType, " a eliminar:")
	fmt.Scanln(&idUser)

	//busca el indice del profesor a eliminar en el slice de usuarios si lo encuentra se asigan el valor
	index := -1
	for i, user := range users {
		if user.id == idUser && user.role == userType {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("el id no corresponde a un", userType)
		fmt.Scanln()
	}

	users = append(users[:index], users[index+1:]...)

	fmt.Println("Profesor eliminado exitosamente.")
	printLine()
	printSimpleText("Presione una tecla para continuar...")
	fmt.Scanln()
	return users
}

func createNewMessage() {
	clearScreen()
	printSimpleText("Crear Mensaje :")
	printLine()
	prompt := promptui.Prompt{
		Label: "Escribe aqui tu mensaje",
		Validate: func(input string) error {
			fmt.Printf("\rLongitud actual: %d", len(input))
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	}

	// Muestra el prompt y captura la entrada del usuario
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	} else {
		fmt.Println("hoaa......")
	}

	printLine()
	printLine()
	fmt.Println("Longitud:", len(result))
	printSimpleText("\nPresione una tecla para continuar")
	fmt.Scanln()
	generateTextFile(result)
}

func generateTextFile(text string) {
	fileName := fmt.Sprintf("mensajes_%s.txt", time.Now().Format("20060102_150405"))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error al crear el archivo: %v\n", err)
		return
	}
	defer file.Close()

	// Escribir los mensajes en el archivo
	_, error := file.WriteString(fmt.Sprintf("%s\n", text))
	if error != nil {
		fmt.Printf("Error al escribir en el archivo: %v\n", error)
		return
	}

	fmt.Printf("Archivo %s creado exitosamente.\n", fileName)
}

// func createAdminRecord() {
// 	fileName := fmt.Sprintf("registro_admin_%s.txt", time.Now().Format("20060102_150405"))
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		fmt.Printf("Error al crear el archivo de registro: %v\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Escribir los mensajes en el archivo
// 	for _, action := range directorActions {
// 		_, err := file.WriteString(fmt.Sprintf("%s: %d\n", action, actionCounter[action]))
// 		if err != nil {
// 			fmt.Printf("Error al escribir en el archivo de registro: %v\n", err)
// 			return
// 		}
// 	}

// 	fmt.Printf("Archivo de registro %s creado exitosamente.\n", fileName)
// }

// func createLaborer() {
// 	laborers = append(laborers, fmt.Sprintf("laborer %d", len(laborers)+1))
// 	fmt.Println("Laborer creado exitosamente.")
// 	for i := 0; i < len(laborers); i++ {
// 		fmt.Println(laborers[i])
// 	}

// }

// func deleteLaborer() {
// 	if len(laborers) == 0 {
// 		fmt.Println("No hay laborers para eliminar.")
// 		return
// 	}

// 	laborers = laborers[:len(laborers)-1]
// 	fmt.Println("Laborer eliminado exitosamente.")
// 	for i := 0; i < len(laborers); i++ {
// 		fmt.Println(laborers[i])
// 	}
// }
