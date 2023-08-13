const table = document.querySelector("table")
const tableBody = document.querySelector("tbody")
const footerCount = document.querySelector(".footer-count")

function addRow({ id, name, age }) {
	const row = document.createElement("tr")
	const idCell = document.createElement("td")
	const nameCell = document.createElement("td")
	const ageCell = document.createElement("td")
	idCell.textContent = id
	nameCell.textContent = name
	ageCell.textContent = age
	row.append(idCell, nameCell, ageCell)
	tableBody.append(row)
	// tableBody.prepend(row)
	footerCount.textContent = tableBody.childElementCount
}

async function fetchPersons() {
	const response = await fetch("/persons")
	const persons = await response.json()
	return persons
}

!async function () {
	const persons = await fetchPersons()
	persons.forEach(addRow)
}()

const webSocketUrl = `ws://${location.host}/persons/live`
console.log(webSocketUrl)
const socket = new WebSocket(webSocketUrl)

socket.addEventListener("open", () => {
	console.log("Connected to websocket")
})

socket.addEventListener("message", event => {
	const person = JSON.parse(event.data)
	addRow(person)
})

socket.addEventListener("close", () => {
	console.log("Disconnected from websocket")
})