import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";
import "./Login.css";
import coffeeImg from "../assets/gambarbiji.png";

function Login() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();

    // 🔹 Simulasi login (nanti bisa diganti dengan Firebase Auth)
    const username = e.target[0].value;
    const password = e.target[1].value;

    if (username === "admin" && password === "1234") {
      alert("Login berhasil sebagai Admin Café!");
      navigate("/admin"); // 🔹 Arahkan ke halaman admin
    } else {
      alert("Username atau password salah!");
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
          <input type="text" placeholder="Username" required />
          <div className="password-container">
            <input
              type={showPassword ? "text" : "password"}
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
          <button type="submit">Login</button>
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
