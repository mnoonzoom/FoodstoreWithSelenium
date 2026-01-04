const API_URL = "http://localhost:8080";
const token = localStorage.getItem("token");
const userId = localStorage.getItem("userId");
console.log("Token:", token);
console.log("User ID:", userId);

if (!token || !userId) {
  window.location.href = "auth.html";
}

const categoryGrids = {
  recommended: document.getElementById("recommendedGrid"),
  appetizers: document.getElementById("appetizersGrid"),
  "main-courses": document.getElementById("mainCoursesGrid"),
  desserts: document.getElementById("dessertsGrid"),
  drinks: document.getElementById("drinksGrid"),
};

let currentPage = 1;
const pageSize = 5;
let totalItems = 0;
let cart = [];

let currentSearch = "";
let currentCategory = "";
let currentSortBy = "price";
let currentSortAsc = true;


// 🍽 Загрузка меню
async function loadMenu() {

  const skip = (currentPage - 1) * pageSize;
  const recommendedBody = {
    limit: pageSize,
    skip,
    search: currentSearch,
    category: "",
    sort_by: currentSortBy,
    sort_asc: currentSortAsc,
  };

  const recommendedRes = await fetch(`${API_URL}/menu/search`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify(recommendedBody),
  });

  if (!recommendedRes.ok) {
    alert("Failed to load recommended menu");
    return;
  }

  const recommendedData = await recommendedRes.json();
  totalItems = recommendedData.total_count;

  Object.values(categoryGrids).forEach(grid => (grid.innerHTML = ""));

  recommendedData.items.forEach(item => {
    renderCard(item, categoryGrids.recommended);
  });

  renderPagination();

  await loadCategories();
}


function renderPagination() {
  const totalPages = Math.ceil(totalItems / pageSize);
  const paginationEl = document.getElementById("recommendedPagination");
  paginationEl.innerHTML = "";

  for (let i = 1; i <= totalPages; i++) {
    const btn = document.createElement("button");
    btn.textContent = i;
    if (i === currentPage) btn.disabled = true;
    btn.addEventListener("click", () => {
      currentPage = i;
      loadMenu();
    });
    paginationEl.appendChild(btn);
  }
}
document.getElementById("filterByNameInput").addEventListener("input", e => {
  currentSearch = e.target.value.trim();
  currentPage = 1;
  loadMenu();
});

document.getElementById("categorySelect").addEventListener("change", e => {
  currentCategory = e.target.value;
  currentPage = 1;
  loadMenu();
});

document.querySelectorAll(".sort-button").forEach(btn => {
  btn.addEventListener("click", () => {
    currentSortBy = "price";
    currentSortAsc = btn.dataset.sort === "asc";
    currentPage = 1;
    loadMenu();
  });
});
function renderCard(item, grid) {
  const card = document.createElement("div");
  card.className = "menu-card";
  card.innerHTML = `
    <div class="img-wrapper">
      <img src="${item.image_url}" alt="${item.name}" onclick='showDishModal(${JSON.stringify(item)})' />
    </div>
    <h3>${item.name}</h3>
    <p>${item.description}</p>
    <strong>$${item.price.toFixed(2)}</strong>
    <button onclick='addToCart(${JSON.stringify(item)})'>Order</button>
  `;
  grid.appendChild(card);
}

async function loadCategories() {
  for (const [category, grid] of Object.entries(categoryGrids)) {
    if (category === "recommended") continue;

    const body = {
      limit: 1000,
      skip: 0,
      search: "",
      category: category,
      sort_by: "price",
      sort_asc: true,
    };

    const res = await fetch(`${API_URL}/menu/search`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify(body),
    });

    if (!res.ok) {
      console.warn(`Failed to load category ${category}`);
      continue;
    }

    const data = await res.json();
    data.items.forEach(item => renderCard(item, grid));
  }
}


function showDishModal(item) {
  document.getElementById("modalDishName").innerText = item.name;
  document.getElementById("modalDishPrice").innerText = item.price.toFixed(2);
  document.getElementById("modalDishDescription").innerText = item.description;
  document.getElementById("modalDishImage").src = item.image_url;
  document.getElementById("addToCartButton").onclick = () => addToCart(item);
  document.getElementById("dishModal").style.display = "flex";
}


function addToCart(item) {
  cart.push(item);
  updateCartUI();
  closeModal("dishModal");
}


function updateCartUI() {
  cartItems.innerHTML = "";
  let total = 0;
  for (const item of cart) {
    total += item.price;
    const div = document.createElement("div");
    div.className = "cart-item";
    div.innerHTML = `<span>${item.name}</span><span>$${item.price.toFixed(
      2
    )}</span>`;
    cartItems.appendChild(div);
  }
  cartTotal.innerText = `$${total.toFixed(2)}`;
}


document.getElementById("checkoutButton").addEventListener("click", () => {
  if (cart.length === 0) return alert("Cart is empty.");
  document.getElementById("checkoutItems").innerHTML = cart
    .map((item) => `<p>${item.name} - $${item.price.toFixed(2)}</p>`)
    .join("");
  document.getElementById("checkoutTotal").innerText = `$${cart
    .reduce((sum, item) => sum + item.price, 0)
    .toFixed(2)}`;
  document.getElementById("checkoutModal").style.display = "flex";
});

document
  .getElementById("confirmOrderButton")
  .addEventListener("click", async () => {
    const itemIds = cart.map((item) => item.id);
    const res = await fetch(`${API_URL}/orders`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_id: userId,
        item_ids: itemIds,
      }),
    });

    const data = await res.json();
    if (res.ok) {
      alert(`Order placed! Order ID: ${data.order_id}`);
      cart = [];
      updateCartUI();
      closeModal("checkoutModal");
    } else {
      alert("Order failed: " + (data.error || data.message));
    }
  });

document
  .querySelectorAll(".modal-close, #modalCloseButton, #closeCheckoutButton")
  .forEach((btn) => {
    btn.addEventListener("click", () => {
      closeModal("dishModal");
      closeModal("checkoutModal");
    });
  });

function closeModal(id) {
  document.getElementById(id).style.display = "none";
}

document.getElementById("logoutButton").addEventListener("click", () => {
  localStorage.clear();
  window.location.href = "auth.html";
});

const authLink = document.getElementById("authLink");
authLink.href = "profile.html";
document.getElementById("authLinkText").innerText = "My Profile";

loadMenu();
