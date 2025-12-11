import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";
import "./Login.css";
import coffeeImg from "./assets/gambarbiji.png";

export default function Login() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [role, setRole] = useState("admin");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const username = e.target.username.value;
    const password = e.target.password.value;

    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, role }),
      });

      // Cek status dulu
      if (!res.ok) {
        const errData = await res.json().catch(() => ({}));
        console.error("Login failed:", errData);
        alert(errData.error || "Login gagal!");
        setLoading(false);
        return;
      }

      // Response valid JSON
      const data = await res.json();

      localStorage.setItem("isLoggedIn", "true");
      localStorage.setItem("role", data.role);
      localStorage.setItem("username", data.username);

      // Navigate sesuai role
      if (data.role === "superadmin" || data.role === "admin") navigate("/adminapprove");
      if (data.role === "cafe") navigate("/admin");
    } catch (err) {
      console.error("Fetch error:", err);
      alert("Server error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login">
      <div className="top-bar-login">
        <h1 className="app-title">
          Carisp
          <img src={coffeeImg} alt="Coffee" className="coffee-img" />
          t
        </h1>
      </div>

      <div className="login-container">
        <h2 className="login-title">LOGIN</h2>
        <form className="login-form" onSubmit={handleSubmit}>
          <input type="text" name="username" placeholder="Username" required />

          <div className="password-container">
            <input
              type={showPassword ? "text" : "password"}
              name="password"
              placeholder="Password"
              required
            />
            <span
              className="toggle-password"
              onClick={() => setShowPassword(!showPassword)}
            >
              {showPassword ? (
                <AiOutlineEyeInvisible size={22} color="#555" />
              ) : (
                <AiOutlineEye size={22} color="#555" />
              )}
            </span>
          </div>

          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="role-dropdown"
            required
          >
            <option value="admin">Admin</option>
            <option value="cafe">Cafe</option>
          </select>

          <button type="submit" disabled={loading}>
            {loading ? "Loading..." : "Login"}
          </button>
        </form>

        <p className="create-account">
          Belum punya akun?{" "}
          <span
            onClick={() => navigate("/register")}
            style={{ cursor: "pointer", color: "#1b2b65", fontWeight: "bold" }}
          >
            Registrasi
          </span>
        </p>
      </div>
    </div>
  );
}
