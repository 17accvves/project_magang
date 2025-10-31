// ProfileCafe.jsx
import React, { useState } from "react";
import {
  FaTrash,
  FaInstagram,
  FaFacebook,
  FaTiktok,
  FaPhone,
  FaClock,
  FaTimes,
  FaPlus,
} from "react-icons/fa";
import "./ProfileCafe.css";

const initialGallery = [
  { id: 1, src: "../../src/assets/LOGO.jpg" },
  { id: 2, src: "../../src/assets/profile.jpg" },
  { id: 3, src: "../../src/assets/maincafe.jpg" },
  { id: 4, src: "../../src/assets/background.png" },
];

function ProfileCafe() {
  const [gallery, setGallery] = useState(initialGallery);
  const [showEdit, setShowEdit] = useState(false);

  const [formData, setFormData] = useState({
    nama: "Yotta",
    alamat: "Jl. Pettrani Raya No.40",
    facebook: "",
    tiktok: "",
    instagram: "",
    telepon: "+62 822-2222-2222",
    jamOperasional: {
      senin: "",
      selasa: "",
      rabu: "",
      kamis: "",
      jumat: "",
      sabtu: "",
      minggu: "",
    },
    fasilitas: {
      wifi: false,
      liveMusic: false,
      outdoor: false,
      airConditioner: false,
    },
    deskripsi: "Lorem ipsum dolor sit amet, consectetur.",
  });

  const handleDeleteImage = (id) => {
    setGallery(gallery.filter((img) => img.id !== id));
  };

  const handleAddImage = (e) => {
    const file = e.target.files[0];
    if (file) {
      const newImage = {
        id: gallery.length + 1,
        src: URL.createObjectURL(file),
      };
      setGallery([...gallery, newImage]);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    if (name in formData.jamOperasional) {
      setFormData({
        ...formData,
        jamOperasional: { ...formData.jamOperasional, [name]: value },
      });
    } else if (name in formData.fasilitas) {
      setFormData({
        ...formData,
        fasilitas: { ...formData.fasilitas, [name]: e.target.checked },
      });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSave = () => {
    console.log("Data tersimpan:", formData);
    setShowEdit(false);
  };

  return (
    <div className="profile-cafe-container">
      <h2 className="profile-title">Profile Cafe</h2>

      {/* Bagian atas */}
      <div className="top-section">
        <div className="main-image-card">
          <img
            src="../../src/assets/maincafe.jpg"
            alt="Main Cafe"
            className="main-image"
          />
          <button className="delete-btn">
            <FaTrash />
          </button>
          <div className="verified">✔ Verified</div>
        </div>

        <div className="info-card">
          <div className="info-row">
            <span className="label">Nama :</span> <span>{formData.nama}</span>
          </div>
          <div className="info-row">
            <span className="label">Alamat :</span>{" "}
            <span>{formData.alamat}</span>
          </div>
          <div className="info-row">
            <span className="label">Sosial Media :</span>
            <span className="social-icons">
              <FaInstagram /> <FaTiktok /> <FaFacebook />
            </span>
          </div>
          <div className="info-row">
            <span className="label">Kontak :</span> <FaPhone />{" "}
            {formData.telepon}
          </div>
          <div className="info-row">
            <span className="label">Jam Operasional :</span> <FaClock /> 08.00 - 23.00
          </div>
          <div className="info-row">
            <span className="label">Fasilitas :</span> ✔ Wifi ✔ Live Music ✔ Outdoor
          </div>
          <div className="info-row">
            <span className="label">Deskripsi :</span> {formData.deskripsi}
          </div>

          <div className="buttons">
            <button className="edit-btn" onClick={() => setShowEdit(true)}>
              Edit Profile ✏️
            </button>
            <button className="save-btn">Simpan Profile 💾</button>
          </div>
        </div>
      </div>

      {/* Galeri */}
      <div className="gallery-section">
        <div className="gallery-card">
          {/* Tombol tambah gambar */}
          <div className="add-image">
            <label htmlFor="fileInput" className="add-image-btn">
              <FaPlus /> Tambah Gambar
            </label>
            <input
              id="fileInput"
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={handleAddImage}
            />
          </div>

          <div className="gallery-grid">
            {gallery.map((img, index) => (
              <div key={img.id} className="gallery-item">
                <img src={img.src} alt={`Cafe ${index + 1}`} />
                <button
                  className="delete-btn"
                  onClick={() => handleDeleteImage(img.id)}
                >
                  <FaTrash />
                </button>
                <span className="pagination">{index + 1}/10</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Popup Edit */}
      {showEdit && (
        <div className="popup-overlay">
          <div className="popup-content">
            <button className="close-btn" onClick={() => setShowEdit(false)}>
              <FaTimes />
            </button>
            <h2>Edit Profile Cafe</h2>

            <label>Nama Cafe:</label>
            <input name="nama" value={formData.nama} onChange={handleInputChange} />

            <label>Alamat:</label>
            <input name="alamat" value={formData.alamat} onChange={handleInputChange} />

            <label>Facebook:</label>
            <input name="facebook" value={formData.facebook} onChange={handleInputChange} />

            <label>TikTok:</label>
            <input name="tiktok" value={formData.tiktok} onChange={handleInputChange} />

            <label>Instagram:</label>
            <input name="instagram" value={formData.instagram} onChange={handleInputChange} />

            <label>Nomor Telepon:</label>
            <input name="telepon" value={formData.telepon} onChange={handleInputChange} />

            <label>Jam Operasional:</label>
            {Object.keys(formData.jamOperasional).map((day) => (
              <div key={day}>
                <span>{day}:</span>
                <input
                  name={day}
                  value={formData.jamOperasional[day]}
                  onChange={handleInputChange}
                />
              </div>
            ))}

            <label>Fasilitas:</label>
            {Object.keys(formData.fasilitas).map((fac) => (
              <div key={fac}>
                <input
                  type="checkbox"
                  name={fac}
                  checked={formData.fasilitas[fac]}
                  onChange={handleInputChange}
                />
                <span>{fac}</span>
              </div>
            ))}

            <label>Deskripsi:</label>
            <textarea
              name="deskripsi"
              value={formData.deskripsi}
              onChange={handleInputChange}
            />

            <div className="popup-buttons">
              <button className="cancel-btn" onClick={() => setShowEdit(false)}>
                Batal ❌
              </button>
              <button className="save-btn" onClick={handleSave}>
                Simpan 💾
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default ProfileCafe;
