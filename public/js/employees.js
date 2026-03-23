// ===============================
// LOAD EMPLOYEES (INDEX PAGE)
// ===============================
async function loadEmployees() {
    const table = document.querySelector("#employeeTable tbody");
    if (!table) return;

    try {
        const response = await fetch("/api/employees");
        const result = await response.json();

        table.innerHTML = "";

        result.data.employees.forEach(emp => {
            const row = document.createElement("tr");

            row.innerHTML = `
                <td><img src="${emp.img_url}" /></td>
                <td>${emp.name}</td>
                <td>${emp.email}</td>
                <td>${emp.designation}</td>
                <td>${emp.department}</td>
                <td>${emp.city}</td>
                <td>${emp.is_active ? "Active" : "Inactive"}</td>
                <td>
                    <button data-id="${emp.id}" class="edit-btn">Edit</button>
                    <button data-id="${emp.id}" class="delete-btn">Delete</button>
                </td>
            `;

            table.appendChild(row);
        });

        attachEditListeners();
        attachDeleteListeners();


    } catch (error) {
        console.error("Failed to load employees:", error);
    }
}

// ===============================
// HANDLE EDIT CLICK (Better than inline onclick)
// ===============================
function attachEditListeners() {
    const buttons = document.querySelectorAll(".edit-btn");

    buttons.forEach(btn => {
        btn.addEventListener("click", () => {
            const id = btn.getAttribute("data-id");
            window.location.href = `add.html?id=${id}`;
        });
    });
}

function attachDeleteListeners() {
    const buttons = document.querySelectorAll(".delete-btn");

    buttons.forEach(btn => {
        btn.addEventListener("click", async () => {
            const id = btn.getAttribute("data-id");

            const confirmDelete = confirm("Are you sure you want to delete this employee?");
            if (!confirmDelete) return;

            try {
                const response = await fetch(`/api/employees/${id}`, {
                    method: "DELETE"
                });

                const result = await response.json();

                if (response.ok) {
                    alert("Employee deleted successfully ✅");

                    // 🔥 Remove row instantly (better UX)
                    btn.closest("tr").remove();

                    // OR reload full list:
                    // loadEmployees();

                } else {
                    alert(result.message || "Delete failed ❌");
                }

            } catch (error) {
                console.error("Delete error:", error);
                alert("Something went wrong ❌");
            }
        });
    });
}

// ===============================
// INIT
// ===============================
document.addEventListener("DOMContentLoaded", () => {
    loadEmployees();
});