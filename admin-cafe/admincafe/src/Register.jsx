import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./Register.css"; 

export default function RegisterCafe() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData();
    formData.append("username", e.target.username.value);
    formData.append("password", e.target.password.value);
    formData.append("email", e.target.email.value);
    formData.append("izin_usaha", e.target.izin_usaha.files[0]);

    try {
      const res = await fetch("http://localhost:8080/register-cafe", {
        method: "POST",
        body: formData,
      });

      const data = await res.json();
      alert(data.message);
      if (res.ok) navigate("/login");
    } catch (err) {
      console.error(err);
      alert("Server error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-container">
        <h2 className="auth-title">Registrasi Cafe</h2>
        <form className="auth-form" onSubmit={handleSubmit}>
          <input type="text" name="username" placeholder="Username" required />
          <input type="email" name="email" placeholder="Email" required />
          <input type="password" name="password" placeholder="Password" required />
          <input type="file" name="izin_usaha" accept=".pdf,.jpg,.png" required />
          <button type="submit">{loading ? "Loading..." : "Registrasi"}</button>
        </form>
        <p className="create-account">
          Sudah punya akun?{" "}
          <span
            onClick={() => navigate("/login")}
            style={{ cursor: "pointer", color: "#1b2b65", fontWeight: "bold" }}
          >
            Login
          </span>
        </p>
      </div>
    </div>
  );
}
