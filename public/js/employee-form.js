const form = document.getElementById("employeeForm");
const message = document.getElementById("message");

const urlParams = new URLSearchParams(window.location.search);
const id = urlParams.get("id");

// ================= INIT =================
document.addEventListener("DOMContentLoaded", async () => {
    await loadFilters();

    if (id) {
        setEditMode();
        await loadEmployeeData(id);
    }

    document.getElementById("addMobileBtn").addEventListener("click", addMobile);
});

// ================= EDIT MODE =================
function setEditMode() {
    document.getElementById("title").innerText = "Edit Employee";
    document.getElementById("submitBtn").innerText = "Update";
}

// ================= LOAD EMPLOYEE =================
async function loadEmployeeData(id) {
    try {
        const res = await fetch(`api/employees/${id}`);
        const result = await res.json();

        const emp = result.data.employee;
        if (!emp) return;

        prefillForm(emp);

    } catch (err) {
        showError("Failed to load employee");
    }
}

// ================= PREFILL =================
function prefillForm(emp) {
    setValue("name", emp.name);
    setValue("designation", emp.designation);
    setValue("department", emp.department);
    setValue("is_active", emp.is_active ? "true" : "false");
    setValue("img_url", emp.img_url);
    setValue("email", emp.email);
    setValue("city", emp.city);
    setValue("country", emp.country);
    setValue("joining_date", emp.joining_date);

    if (emp.mobiles && emp.mobiles.length > 0) {
        const container = document.getElementById("mobileContainer");
        container.innerHTML = "";

        emp.mobiles.forEach(m => addMobile(m.type, m.number));
    }
}

// ================= MOBILE =================
function addMobile(type = "home", number = "") {
    const container = document.getElementById("mobileContainer");

    if (container.children.length >= 3) {
        showError("Max 3 mobiles allowed");
        return;
    }

    const div = document.createElement("div");
    div.className = "mobile-row";

    div.innerHTML = `
        <select class="mobile-type">
            <option value="home">Home</option>
            <option value="office">Office</option>
            <option value="other">Other</option>
        </select>
        <input type="text" class="mobile-number" placeholder="Mobile number" />
    `;

    container.appendChild(div);

    div.querySelector(".mobile-type").value = type;
    div.querySelector(".mobile-number").value = number;
}

function getMobiles() {
    const rows = document.querySelectorAll(".mobile-row");
    const mobiles = [];

    rows.forEach(row => {
        const type = row.querySelector(".mobile-type").value;
        const number = row.querySelector(".mobile-number").value.trim();

        if (number) {
            mobiles.push({ type, number });
        }
    });

    return mobiles;
}

// ================= FILTERS =================
async function loadFilters() {
    try {
        const res = await fetch("api/employees/filters");
        const result = await res.json();

        populateDropdown("designation", result.data.designations);
        populateDropdown("department", result.data.departments);

    } catch {
        console.error("Filter load error");
    }
}

function populateDropdown(id, items) {
    const el = document.getElementById(id);
    if (!el) return;

    el.innerHTML = `<option value="">Select ${id}</option>`;
    items.forEach(i => {
        const opt = document.createElement("option");
        opt.value = i;
        opt.textContent = i;
        el.appendChild(opt);
    });
}

// ================= SUBMIT =================
if (form) {
    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const data = {
            name: getValue("name"),
            designation: getValue("designation"),
            department: getValue("department"),
            is_active: getValue("is_active") === "true",
            img_url: getValue("img_url"),
            email: getValue("email"),
            city: getValue("city"),
            country: getValue("country"),
            joining_date: getValue("joining_date") || null,
            mobiles: getMobiles()
        };

        try {
            let res;

            if (id) {
                res = await fetch(`/api/employees/${id}`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });
            } else {
                res = await fetch("/api/employees/", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });
            }

            if (res.ok) {
                window.location.href = "index.html";
            } else {
                showError(await res.text());
            }

        } catch {
            showError("Network error");
        }
    });
}

// ================= HELPERS =================
function getValue(id) {
    return document.getElementById(id)?.value || "";
}

function setValue(id, value) {
    const el = document.getElementById(id);
    if (!el) return;
    el.value = value || "";
}

function showError(msg) {
    message.innerText = "❌ " + msg;
}