
let allOrdersAdmin = [];
let allMenuItems = [];

function loadAdminPanel(token, userId) {
    const adminPanel = document.getElementById("adminPanel");
    if (!adminPanel) return;
    adminPanel.style.display = "block";

    window.loadOrders = async function () {
        try {
            const res = await fetch(`${API_URL}/orders`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!res.ok) throw new Error("Failed to load orders");

            allOrdersAdmin = await res.json();

            const uniqueItemIds = [...new Set(allOrdersAdmin.flatMap(order => order.item_ids))];
            await loadMenuItems(uniqueItemIds);

            currentPage = 1;
            renderOrdersAdminPage();
        } catch (err) {
            console.error(err);
            alert("Failed to load orders.");
        }
    };

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

        allMenuItems = data;
        menuMap = {};
        data.forEach(item => {
            menuMap[item.id] = item.name;
        });
    }

    function renderOrdersAdminPage() {
        const container = document.getElementById("adminContent");
        if (!container) return;

        container.innerHTML = "<h3>Manage Orders</h3>";

        const start = (currentPage - 1) * pageSize;
        const end = start + pageSize;
        const pageOrders = allOrdersAdmin.slice(start, end);

        if (pageOrders.length === 0) {
            container.innerHTML += "<p>No orders found.</p>";
            return;
        }

        pageOrders.forEach(order => {
            const itemNames = order.item_ids.map(id => menuMap[String(id).trim()] || id).join(", ");
            const div = document.createElement("div");
            div.classList.add("order-entry");
            div.innerHTML = `
        <p><strong>Order ID:</strong> ${order.id}</p>
        <p><strong>Status:</strong> ${order.status}</p>
        <p><strong>Total:</strong> $${Number(order.total_price || 0).toFixed(2)}</p>
        <p><strong>Items:</strong> ${itemNames}</p>
        <p><strong>Date:</strong> ${new Date(order.created_at).toLocaleString()}</p>
        <button onclick="updateOrderStatus('${order.id}')">Update Status</button>
        <button onclick="deleteOrder('${order.id}')">Delete Order</button>
        <hr>
      `;
            container.appendChild(div);
        });

        container.innerHTML += `
      <div class="pagination-controls">
        <button onclick="prevPageOrders()" class="btn btn-primary">Previous</button>
        <span>Page ${currentPage}</span>
        <button onclick="nextPageOrders()" class="btn btn-primary">Next</button>
      </div>
    `;
    }

    window.prevPageOrders = () => {
        if (currentPage > 1) {
            currentPage--;
            renderOrdersAdminPage();
        }
    };

    window.nextPageOrders = () => {
        if (currentPage * pageSize < allOrdersAdmin.length) {
            currentPage++;
            renderOrdersAdminPage();
        }
    };

    window.deleteOrder = async (orderId) => {
        if (!confirm("Are you sure you want to delete this order?")) return;
        try {
            const res = await fetch(`${API_URL}/orders/${orderId}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!res.ok) throw new Error("Failed to delete order");
            alert("Order deleted");
            await window.loadOrders();
        } catch (err) {
            console.error(err);
            alert("Error deleting order");
        }
    };

    window.updateOrderStatus = async (orderId) => {
        const newStatus = prompt("Enter new status for order:", "Pending");
        if (!newStatus) return;
        try {
            const res = await fetch(`${API_URL}/orders/${orderId}/status`, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({ status: newStatus }),
            });
            if (!res.ok) throw new Error("Failed to update order");
            alert("Order updated");
            await window.loadOrders();
        } catch (err) {
            console.error(err);
            alert("Error updating order");
        }
    };


    window.loadMenu = async function () {
        try {
            const res = await fetch(`${API_URL}/menu`, {
                headers: { Authorization: `Bearer ${token}` },  cache: "no-cache"
            });
            if (!res.ok) throw new Error("Failed to load menu");
            allMenuItems = await res.json();
            renderMenu();
        } catch (err) {
            console.error(err);
            alert("Failed to load menu");
        }
    };

    function renderMenu() {
        const container = document.getElementById("adminContent");
        container.innerHTML = "<h3>Manage Menu</h3>";

        allMenuItems.forEach(item => {
            const div = document.createElement("div");
            div.classList.add("menu-entry");
            div.innerHTML = `
        <p><strong>${item.name}</strong> - $${Number(item.price).toFixed(2)}</p>
        <p>Category: ${item.category}</p>
        <button onclick="showEditMenuForm('${item.id}')">Edit</button>
        <button onclick="deleteMenuItem('${item.id}')">Delete</button>
        <hr>
      `;
            container.appendChild(div);
        });
    }

    window.deleteMenuItem = async (itemId) => {
        if (!confirm("Delete this menu item?")) return;
        try {
            const res = await fetch(`${API_URL}/menu/${itemId}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!res.ok) throw new Error("Failed to delete menu item");
            alert("Menu item deleted");
            await window.loadMenu();
        } catch (err) {
            console.error(err);
            alert("Error deleting menu item");
        }
    };

    window.showAddMenuForm = function () {
        document.getElementById("addMenuForm").style.display = "block";
    };

    window.hideAddMenuForm = function () {
        document.getElementById("addMenuForm").style.display = "none";
    };

    document.getElementById("menuForm").addEventListener("submit", async (e) => {
        e.preventDefault();

        const newItem = {
            name: document.getElementById("menuItemName").value,
            price: parseFloat(document.getElementById("menuItemPrice").value),
            category: document.getElementById("menuItemCategory").value,
            image_url: document.getElementById("menuItemPicture").value,
            available: document.getElementById("menuItemAvailable").checked,

        };

        try {
            const res = await fetch(`${API_URL}/menu`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(newItem),
            });

            if (!res.ok) throw new Error("Failed to add menu item");
            alert("Menu item added");
            hideAddMenuForm();
            await window.loadMenu();
        } catch (err) {
            console.error(err);
            alert("Error adding menu item");
        }
    });
    let currentEditItemId = null;

    window.showEditMenuForm = function(itemId) {
        currentEditItemId = itemId;
        const item = allMenuItems.find(i => i.id === itemId);
        if (!item) {
            alert("Menu item not found");
            return;
        }

        const modal = document.getElementById("editMenuModal");
        if (!modal) {
            alert("Edit modal not found");
            return;
        }

        modal.style.display = "flex";  // или "block"

        const form = document.getElementById("editMenuForm");
        form.querySelector("#editMenuItemName").value = item.name;
        form.querySelector("#editMenuItemPrice").value = item.price;
        form.querySelector("#editMenuItemCategory").value = item.category;
        form.querySelector("#editMenuItemPicture").value = item.image_url;
        form.querySelector("#editMenuItemAvailable").checked = item.available || false;

    }

    window.hideEditMenuForm = function() {
        const modal = document.getElementById("editMenuModal");
        if (modal) modal.style.display = "none";
        currentEditItemId = null;
    }

    document.getElementById("editMenuForm").addEventListener("submit", async (e) => {
        e.preventDefault();

        if (!currentEditItemId) {
            alert("No menu item selected for editing");
            return;
        }

        const updatedItem = {
            id: currentEditItemId,
            name: document.getElementById("editMenuItemName").value,
            price: parseFloat(document.getElementById("editMenuItemPrice").value),
            category: document.getElementById("editMenuItemCategory").value,
            image_url: document.getElementById("editMenuItemPicture").value,
        };

        try {
            const res = await fetch(`${API_URL}/menu/${currentEditItemId}`, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`,
                },
                body: JSON.stringify(updatedItem),
            });

            if (!res.ok) throw new Error("Failed to update menu item");

            alert("Menu item updated");
            hideEditMenuForm();
            await loadMenu();
        } catch (err) {
            console.error(err);
            alert("Error updating menu item");
        }
    });

}
