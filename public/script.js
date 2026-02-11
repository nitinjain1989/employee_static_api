// Load employees (only runs if table exists)
async function loadEmployees() {
    const table = document.querySelector("#employeeTable tbody");
    if (!table) return;

    const response = await fetch("/employees");
    const result = await response.json();

    table.innerHTML = "";

    result.data.employees.forEach(emp => {
        table.innerHTML += `
      <tr>
        <td><img src="${emp.img_url}" /></td>
        <td>${emp.name}</td>
        <td>${emp.email}</td>
        <td>${emp.designation}</td>
        <td>${emp.department}</td>
        <td>${emp.city}</td>
        <td>${emp.is_active ? "Active" : "Inactive"}</td>
        <td>
            <button onclick="editEmployee('${emp.id}')">Edit</button>
         </td>
      </tr>
    `;
    });
}

// ===============================
// NAVIGATE TO EDIT PAGE
// ===============================
function editEmployee(id) {
    window.location.href = `add.html?id=${id}`;
}


// -------- ADD + EDIT EMPLOYEE (same page) --------

const form = document.getElementById("employeeForm");
const message = document.getElementById("message");

// Get ID from URL
const urlParams = new URLSearchParams(window.location.search);
const id = urlParams.get("id");

// If ID exists → Edit Mode
if (id) {
    document.querySelector("h2").innerText = "Edit Employee";
    loadEmployeeData(id);
}

async function loadEmployeeData(id) {
    try {
        const response = await fetch("/employees");
        const result = await response.json();

        const employee = result.data.employees.find(e => e.id === id);
        if (!employee) return;

        document.getElementById("name").value = employee.name || "";
        document.getElementById("designation").value = employee.designation || "";
        document.getElementById("department").value = employee.department || "";
        document.getElementById("is_active").value = employee.is_active ? "true" : "false";
        document.getElementById("img_url").value = employee.img_url || "";
        document.getElementById("email").value = employee.email || "";
        document.getElementById("city").value = employee.city || "";
        document.getElementById("country").value = employee.country || "";
        document.getElementById("joining_date").value = employee.joining_date || "";

    } catch (err) {
        console.error("Error loading employee:", err);
    }
}

if (form) {
    form.addEventListener("submit", async function (e) {
        e.preventDefault();

        const dateValue = document.getElementById("joining_date").value;

        const data = {
            name: document.getElementById("name").value,
            designation: document.getElementById("designation").value,
            department: document.getElementById("department").value,
            is_active: document.getElementById("is_active").value === "true",
            img_url: document.getElementById("img_url").value,
            email: document.getElementById("email").value,
            city: document.getElementById("city").value,
            country: document.getElementById("country").value,
            joining_date: dateValue ? dateValue : null
        };

        try {
            let response;

            if (id) {
                // UPDATE
                response = await fetch(`/employees/update?id=${id}`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });
            } else {
                // CREATE
                response = await fetch("/employees/create", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });
            }

            if (response.ok) {
                window.location.href = "index.html";
            } else {
                const text = await response.text();
                message.innerText = "❌ Error: " + text;
                message.style.color = "red";
            }

        } catch (err) {
            message.innerText = "❌ Network error";
            message.style.color = "red";
        }
    });
}

loadEmployees();
