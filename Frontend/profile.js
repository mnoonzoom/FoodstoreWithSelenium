let allOrders = [];
let menuMap = {};
let currentPage1 = 1;
const pageSize1 = 4;

document.addEventListener("DOMContentLoaded", () => {
  const token = localStorage.getItem("token");
  const userId = localStorage.getItem("userId");

  if (!token || !userId) {
    alert("You must be logged in.");
    window.location.href = "auth.html";
    return;
  }

  const logoutBtn = document.getElementById("logoutButton");
  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.clear();
      window.location.href = "auth.html";
    });
  }



  async function loadProfile() {
    try {
      const res = await fetch(`${API_URL}/users/${userId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Profile error");

      const profileName = document.getElementById("profileName");
      if (profileName) profileName.innerText = `Name: ${data.username || "N/A"}`;

      const profileEmail = document.getElementById("profileEmail");
      if (profileEmail) profileEmail.innerText = `Email: ${data.email || "N/A"}`;

      const profilePhone = document.getElementById("profilePhone");
      if (profilePhone) profilePhone.innerText = `Phone: ${data.phone || "N/A"}`;

      if (data.role === "admin") {
        if (typeof loadAdminPanel === "function") {
          loadAdminPanel(token, userId);
        }
      }
    } catch (err) {
      console.error(err);
      alert("Failed to load profile.");
    }
  }

  async function loadAllOrders() {
    try {
      const res = await fetch(`${API_URL}/orders/user/${userId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error("Failed to load orders");

      allOrders = await res.json();

      const uniqueItemIds = [...new Set(allOrders.flatMap(order => order.item_ids))];
      await loadMenuItems(uniqueItemIds);

      currentPage1 = 1;
      renderOrdersPage();
    } catch (err) {
      console.error(err);
    
    }
  }

  async function loadMenuItems(ids) {
    if (ids.length === 0) return;

    const res = await fetch(`${API_URL}/menu/multiple`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({ ids }),
    });

    if (!res.ok) {
      throw new Error("Failed to load menu items");
    }

    const data = await res.json();
    menuMap = {};
    data.forEach(item => {
      menuMap[item.id] = item.name;
    });
  }

  function renderOrdersPage() {
    const container = document.getElementById("orderHistoryContainer");
    if (!container) return;

    container.innerHTML = "";

    const start = (currentPage1 - 1) * pageSize1;
    const end = start + pageSize1;
    const pageOrders = allOrders.slice(start, end);

    if (pageOrders.length === 0) {
      container.innerHTML = "<p>No orders found.</p>";
      return;
    }

    pageOrders.forEach((order) => {
      const itemNames = order.item_ids.map(id => menuMap[String(id).trim()] || id).join(", ");

      const div = document.createElement("div");
      div.classList.add("order-entry");
      div.innerHTML = `
        <p><strong>Order ID:</strong> ${order.id}</p>
        <p><strong>Status:</strong> ${order.status}</p>
        <p><strong>Total:</strong> $${Number(order.total_price || 0).toFixed(2)}</p>
        <p><strong>Items:</strong> ${itemNames}</p>
        <p><strong>Date:</strong> ${new Date(order.created_at).toLocaleString()}</p>
        <hr>
      `;
      container.appendChild(div);
    });

    const pageLabel = document.getElementById("currentPageLabel");
    if (pageLabel) pageLabel.innerText = `Page ${currentPage1}`;
  }

  document.getElementById("prevPageButton").addEventListener("click", () => {
    if (currentPage1 > 1) {
      currentPage1--;
      renderOrdersPage();
    }
  });

  document.getElementById("nextPageButton").addEventListener("click", () => {
    if (currentPage1 * pageSize1 < allOrders.length) {
      currentPage1++;
      renderOrdersPage();
    }
  });

  loadProfile();
  loadAllOrders();
});
