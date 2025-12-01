import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";
import "./Login.css";
import coffeeImg from "./assets/gambarbiji.png";

function Login() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [role, setRole] = useState("admin");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const username = e.target.username.value;
    const password = e.target.password.value;

    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password, role }),
      });

      const data = await res.json();

      if (!res.ok) {
        alert(data.error || "Login gagal!");
      } else {
        alert(`Login berhasil sebagai ${role.toUpperCase()}!`);
        // Redirect sesuai role
        if (role === "admin") navigate("/superadmin");
        else if (role === "cafe") navigate("/admin");
        else if (role === "eo") navigate("/eo");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan server!");
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
          <input
            type="text"
            name="username"
            placeholder="Username"
            required
          />

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
            <option value="eo">EO</option>
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

export default Login;
