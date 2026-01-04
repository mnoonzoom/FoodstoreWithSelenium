
const API_URL = "http://localhost:8080";

document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();

    const email = document.getElementById("loginEmail").value;
    const password = document.getElementById("loginPassword").value;
    console.log(email, password); // ← вот оно
  const res = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const data = await res.json();

  if (res.ok && data.token) {
    localStorage.setItem("token", data.token);
    localStorage.setItem("userId", data.user_id);
    window.location.href = "index.html";

  } else {
    alert("Login failed: " + (data.error || data.message));
  }
    console.log("Login response data:", data);
});
document.querySelectorAll('.auth-tab').forEach(tab => {
    tab.addEventListener('click', () => {
        const target = tab.getAttribute('data-tab');

        document.querySelectorAll('.auth-tab').forEach(t => t.classList.remove('active'));
        tab.classList.add('active');

        const loginForm = document.getElementById('loginForm');
        const registerForm = document.getElementById('registerForm');

        if (target === 'login') {
            loginForm.classList.remove('hidden');
            registerForm.classList.add('hidden');
        } else if (target === 'register') {
            registerForm.classList.remove('hidden');
            loginForm.classList.add('hidden');
        }
    });
});

document
    .getElementById("registerForm")
    .addEventListener("submit", async (e) => {
        e.preventDefault();

        const name = document.getElementById("registerName").value;
        const email = document.getElementById("registerEmail").value;
        const password = document.getElementById("registerPassword").value;
        const confirmPassword = document.getElementById(
            "registerConfirmPassword"
        ).value;

        if (password !== confirmPassword) {
            return alert("Passwords do not match.");
        }

        const res = await fetch(`${API_URL}/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username: name, email, password }),
        });

        const data = await res.json();

        if (res.ok) {
            alert("Registration successful! Please log in.");
            document.querySelector('[data-tab="login"]').click();
        } else {
            alert("Registration failed: " + (data.error || data.message));
        }
    });
