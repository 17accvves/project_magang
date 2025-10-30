import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";
import "./Register.css";
import coffeeImg from "./assets/gambarbiji.png";

function Register() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [showPopup, setShowPopup] = useState(false); // 🔹 untuk menampilkan popup

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
    if (selectedFile && selectedFile.type.startsWith("image/")) {
      setPreview(URL.createObjectURL(selectedFile));
    } else {
      setPreview(null);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!file) {
      alert("Harap unggah surat izin usaha terlebih dahulu!");
      return;
    }

    // Simulasi kirim data registrasi
    console.log("Data registrasi dikirim...");

    // Tampilkan popup "menunggu konfirmasi admin"
    setShowPopup(true);
  };

  const handleClosePopup = () => {
    setShowPopup(false);
    // Simulasi bahwa admin sudah menyetujui pendaftaran (di dunia nyata: ditunggu dari backend)
    setTimeout(() => {
      alert("Akun telah disetujui oleh admin!");
      navigate("/login");
    }, 5000);
  };

  return (
    <div className="Register-page">
      <div className="top-bar-register">
        <h1 className="app-title">
          Carisp
          <img src={coffeeImg} alt="Coffee" className="coffee-img" />
          t
        </h1>
      </div>

      <div className="Register-container">
        <h2 className="register-title">REGISTRASI AKUN</h2>

        <form className="register-form" onSubmit={handleSubmit}>
          <input type="text" placeholder="Nama Kafe" required />
          <input type="email" placeholder="Email" required />

          <div className="form-group">
            <label className="title-usaha">Surat Izin Usaha</label>
            <label className="custom-file-upload">
              <input
                type="file"
                accept=".jpg,.jpeg,.png,.pdf"
                onChange={handleFileChange}
              />
            </label>
            {file && <p className="file-name">📄 {file.name}</p>}
          </div>

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

          <button type="submit">Registrasi</button>
        </form>

        <p className="sudah-punya-akun">
          <span
            onClick={() => navigate("/login")}
            style={{ cursor: "pointer", color: "#1b2b65", fontWeight: "bold" }}
          >
            Sudah punya akun? Login
          </span>
        </p>
      </div>

      {/* 🔹 Popup Konfirmasi */}
      {showPopup && (
        <div className="popup-overlay">
          <div className="popup">
            <h3>Menunggu Konfirmasi Admin</h3>
            <p>
              Data registrasi kamu sudah terkirim. Mohon tunggu konfirmasi dari
              admin sebelum bisa login.
            </p>
            <button onClick={handleClosePopup}>OK</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default Register;
