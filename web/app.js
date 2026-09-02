const form = document.getElementById("application-form");
const tbody = document.getElementById("applications-body");
const statusFilter = document.getElementById("status-filter");
const formMessage = document.getElementById("form-message");

const statusLabels = {
    new: "Новая",
    in_progress: "В работе",
    success: "Успешно",
    rejected: "Отказ"
};

async function loadApplications() {
    const status = statusFilter.value;

    let url = "/admin/applications";

    if (status) {
        url += `?status=${encodeURIComponent(status)}`;
    }

    const response = await fetch(url);

    if (!response.ok) {
        tbody.innerHTML = `
            <tr>
                <td colspan="6">Не удалось загрузить заявки</td>
            </tr>
        `;
        return;
    }

    const data = await response.json();

    renderApplications(data.items || []);
}

function renderApplications(applications) {
    tbody.innerHTML = "";

    if (applications.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="6">Заявок нет</td>
            </tr>
        `;
        return;
    }

    for (const application of applications) {
        const row = document.createElement("tr");

        row.innerHTML = `
            <td>${escapeHtml(application.name)}</td>
            <td>${escapeHtml(application.phone)}</td>
            <td>${escapeHtml(application.source)}</td>
            <td>${escapeHtml(application.comment)}</td>
            <td>
                <select class="status-select" data-id="${application.id}">
                    ${renderStatusOptions(application.status)}
                </select>
            </td>
            <td>${formatDate(application.created_at)}</td>
        `;

        tbody.appendChild(row);
    }

    document.querySelectorAll(".status-select").forEach(select => {
        select.addEventListener("change", updateStatus);
    });
}

function renderStatusOptions(currentStatus) {
    return Object.entries(statusLabels)
        .map(([value, label]) => {
            const selected = value === currentStatus ? "selected" : "";
            return `<option value="${value}" ${selected}>${label}</option>`;
        })
        .join("");
}

async function updateStatus(event) {
    const select = event.target;
    const id = select.dataset.id;
    const status = select.value;

    const response = await fetch(`/admin/applications/${id}/status`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ status })
    });

    if (!response.ok) {
        alert("Не удалось изменить статус");
        await loadApplications();
        return;
    }

    await loadApplications();
}

form.addEventListener("submit", async event => {
    event.preventDefault();

    formMessage.textContent = "";
    formMessage.className = "";

    const payload = {
        name: document.getElementById("name").value,
        phone: document.getElementById("phone").value,
        comment: document.getElementById("comment").value,
        source: document.getElementById("source").value
    };

    const response = await fetch("/admin/applications", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    });

    if (!response.ok) {
        formMessage.textContent = "Не удалось создать заявку";
        formMessage.className = "error";
        return;
    }

    form.reset();

    formMessage.textContent = "Заявка создана";
    formMessage.className = "success";

    await loadApplications();
});

statusFilter.addEventListener("change", loadApplications);

function formatDate(value) {
    return new Date(value).toLocaleString("ru-RU");
}

function escapeHtml(value) {
    if (!value) {
        return "";
    }

    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#039;");
}

loadApplications();