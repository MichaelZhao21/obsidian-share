// Global data object
let data = [];

// Cookie helper functions
function setCookie(name, value, days) {
    let expires = "";
    if (days) {
        const date = new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie =
        name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === " ") c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0)
            return c.substring(nameEQ.length, c.length);
    }
    return null;
}

function eraseCookie(name) {
    document.cookie =
        name +
        "=;expires=Thu, 01 Jan 1970 00:00:01 GMT" +
        ";path=/";
}

// Check login status when page loads
document.addEventListener("DOMContentLoaded", checkLoginStatus);

// Login form submission
document
    .getElementById("loginForm")
    .addEventListener("submit", function (e) {
        e.preventDefault();
        login();
    });

// Logout button
document
    .getElementById("logoutButton")
    .addEventListener("click", logout);

// Function to check if user is logged in
async function checkLoginStatus() {
    const res = await fetch("/admin/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ password }),
    });

    if (res.status === 200) {
        // User is logged in
        showLoggedInView();
        fetchUserData();
    } else {
        // User is not logged in
        showLoginView();
    }
}

// Function to handle login
async function login() {
    const password = document.getElementById("password").value;

    // Show loading, hide messages
    document
        .getElementById("loginLoading")
        .classList.remove("hidden");
    document.getElementById("loginError").classList.add("hidden");
    document.getElementById("loginSuccess").classList.add("hidden");

    setCookie("admin_pass", password, 365); // Cookie expires in 1 year

    // POST login
    const res = await fetch("/admin/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (res.status !== 200) {
        // Login failed
        document.getElementById("loginError").textContent =
            "Invalid password";
        document
            .getElementById("loginError")
            .classList.remove("hidden");
        document
            .getElementById("loginLoading")
            .classList.add("hidden");
        return;
    }

    // Login successful
    document.getElementById("loginSuccess").textContent =
        "Login successful! Redirecting...";
    document
        .getElementById("loginSuccess")
        .classList.remove("hidden");

    // Reload the page after successful login
    setTimeout(() => {
        window.location.reload();
    }, 1000);
}

// Function to handle logout
function logout() {
    eraseCookie("admin_pass");
    window.location.reload();
}

// Function to show logged-in view
function showLoggedInView() {
    document
        .getElementById("loginContainer")
        .classList.add("hidden");
    document
        .getElementById("dataContainer")
        .classList.remove("hidden");
}

// Function to show login view
function showLoginView() {
    document
        .getElementById("loginContainer")
        .classList.remove("hidden");
    document
        .getElementById("dataContainer")
        .classList.add("hidden");
}

// Function to fetch user data
async function fetchUserData() {
    const userDataContainer = document.getElementById("userData");

    const res = await fetch("/admin/list", {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });
    if (res.status !== 200) {
        console.error("Failed to fetch user data");
        alert("Failed to fetch user data");
        return;
    }

    data = await res.json();
    console.log(data);

    const dataStr = data
        .map(
            (e, i) => `
                    <div class="data-item">
                        <p id="data-slug-${i}" contenteditable>${e.slug}</p>
                        <p id="data-filename-${i}" contenteditable>${e.filename}</p>
                        <div class="dummy"></div>
                        <button id="data-save-${i}" class="data-button" onclick="save(${i})">Save</button>
                        <button id="data-delete-${i}" class="data-button delete" onclick="deleteItem(${i})">Delete</button>
                    </div>
                `
        )
        .join("\n");

    userDataContainer.innerHTML = `
        <div class="data-items">
            <div id="data-header-buttons" class="data-item">
                <button id="data-add" class="data-button" onclick="add()">Add New...</button>
                <button id="data-pull" class="data-button" onclick="pullMd()">Pull and compile files</button>
            </div>
            <div class="data-item">
                <p><b>Slug</b></p>
                <p><b>Filename</b></p>
            </div>
            ${dataStr}
        </div>
    `;
}

function showMessageAndReturnError(res, successMsg, failureMsg) {
    if (res.status === 200) {
        document.getElementById("data-message").textContent = successMsg;
        document.getElementById("data-message").classList.add("success");
    } else {
        document.getElementById("data-message").textContent = failureMsg;
        document.getElementById("data-message").classList.add("error");
    }

    setTimeout(() => {
        document.getElementById("data-message").textContent = "";
        document.getElementById("data-message").classList.remove("error");
        document.getElementById("data-message").classList.remove("success");
    }, 3000);

    return res.status !== 200;
}

async function save(i) {
    const slug = document.getElementById(`data-slug-${i}`).textContent;
    const filename = document.getElementById(`data-filename-${i}`).textContent;

    const res = await fetch("/admin/data", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ slug, filename }),
    });

    if (showMessageAndReturnError(res, "Data saved", "Error saving data")) {
        return;
    }

    data[i].slug = slug;
    data[i].filename = filename;
}

async function deleteItem(i) {
    const res = await fetch("/admin/data", {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ id: data[i].id }),
    });

    if (showMessageAndReturnError(res, "Data deleted", "Error deleting data")) {
        return;
    }

    // Reload the page on delete
    window.location.reload();
}

function add() {
    data.push({ slug: "", filename: "" });

    const userDataContainer = document.getElementById("userData");

    const newItem = `
        <div class="data-item">
            <p id="data-slug-${data.length - 1}" contenteditable>slug</p>
            <p id="data-filename-${data.length - 1}" contenteditable>filename</p>
            <div class="dummy"></div>
            <button id="data-save-${data.length - 1}" class="data-button" onclick="save(${data.length - 1})">Save</button>
            <button id="data-delete-${data.length - 1}" class="data-button delete" onclick="deleteItem(${data.length - 1})">Delete</button>
        </div>
    `;
    userDataContainer.insertAdjacentHTML("beforeend", newItem);

    document.getElementById(`data-slug-${data.length - 1}`).focus();
}

async function pullMd() {
    document.getElementById("data-message").textContent = "Pulling files...";
    const res = await fetch("/admin/pull", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
    });

    if (showMessageAndReturnError(res, "Files pulled and compiled", "Error pulling files")) {
        return;
    }
}