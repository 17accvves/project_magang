// ProfileCafe.jsx
import React, { useState, useEffect } from "react";
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
  const [statusBuka, setStatusBuka] = useState("");

  const [formData, setFormData] = useState({
    nama: "Yotta",
    alamat: "Jl. Pettarani Raya No.40",
    facebook: "https://facebook.com/",
    tiktok: "https://www.tiktok.com/@",
    instagram: "https://www.instagram.com/",
    telepon: "+62 822-2222-2222",
    jamOperasional: {
      senin: { buka: "08:00", tutup: "23:00" },
      selasa: { buka: "08:00", tutup: "23:00" },
      rabu: { buka: "08:00", tutup: "23:00" },
      kamis: { buka: "08:00", tutup: "23:00" },
      jumat: { buka: "08:00", tutup: "23:00" },
      sabtu: { buka: "08:00", tutup: "23:00" },
      minggu: { buka: "08:00", tutup: "23:00" },
    },
    fasilitas: {
      wifi: true,
      liveMusic: true,
      outdoor: true,
      airConditioner: false,
    },
    deskripsi: "Lorem ipsum dolor sit amet, consectetur.",
  });

  // === LOGIKA CEK HARI DAN JAM BUKA ===
  useEffect(() => {
    const updateStatus = () => {
      const hariList = [
        "minggu",
        "senin",
        "selasa",
        "rabu",
        "kamis",
        "jumat",
        "sabtu",
      ];

      const now = new Date();
      const dayIndex = now.getDay();
      const currentDay = hariList[dayIndex];
      const currentTime = now.getHours() * 60 + now.getMinutes();

      const jadwal = formData.jamOperasional[currentDay];
      const [bukaH, bukaM] = jadwal.buka.split(":").map(Number);
      const [tutupH, tutupM] = jadwal.tutup.split(":").map(Number);

      const bukaMinutes = bukaH * 60 + bukaM;
      const tutupMinutes = tutupH * 60 + tutupM;

      const isOpen =
        currentTime >= bukaMinutes && currentTime <= tutupMinutes;

      setStatusBuka({
        hari: currentDay.charAt(0).toUpperCase() + currentDay.slice(1),
        buka: jadwal.buka,
        tutup: jadwal.tutup,
        status: isOpen ? "Buka Sekarang" : "Tutup Sekarang",
        warna: isOpen ? "#00c896" : "#ff7878",
      });
    };

    updateStatus();
    const interval = setInterval(updateStatus, 60000);
    return () => clearInterval(interval);
  }, [formData]);

  // === Fungsi lainnya ===
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
    const { name, value, type, checked, dataset } = e.target;

    if (dataset.day && dataset.field) {
      setFormData({
        ...formData,
        jamOperasional: {
          ...formData.jamOperasional,
          [dataset.day]: {
            ...formData.jamOperasional[dataset.day],
            [dataset.field]: value,
          },
        },
      });
    } else if (name in formData.fasilitas) {
      setFormData({
        ...formData,
        fasilitas: { ...formData.fasilitas, [name]: checked },
      });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSave = () => {
    console.log("Data tersimpan:", formData);
    setShowEdit(false);
  };

  // === Buka link sosial media ===
  const openSocialLink = (platform) => {
    const link = formData[platform];
    if (!link || link.trim() === "" || link === "https://") {
      alert(`Link ${platform} belum diatur!`);
      return;
    }

    // pastikan ada protokol
    const finalLink = link.startsWith("http")
      ? link
      : `https://${link}`;
    window.open(finalLink, "_blank");
  };

  return (
    <div className="profile-cafe-container">
      <h2 className="profile-title">Profile Cafe</h2>

      {/* Bagian Atas */}
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

          {/* === Bagian Sosial Media === */}
          <div className="info-row">
            <span className="label">Sosial Media :</span>
            <span className="social-icons">
              <FaInstagram
                className="social-icon"
                onClick={() => openSocialLink("instagram")}
              />
              <FaTiktok
                className="social-icon"
                onClick={() => openSocialLink("tiktok")}
              />
              <FaFacebook
                className="social-icon"
                onClick={() => openSocialLink("facebook")}
              />
            </span>
          </div>

          <div className="info-row">
            <span className="label">Kontak :</span> <FaPhone />{" "}
            {formData.telepon}
          </div>

          {/* Menampilkan Hari dan Status Buka */}
          <div className="info-row">
            <span className="label">Jam Operasional ({statusBuka.hari}) :</span>{" "}
            <FaClock /> {statusBuka.buka} - {statusBuka.tutup}
            <span
              style={{
                marginLeft: "10px",
                color: "white",
                background: statusBuka.warna,
                padding: "3px 8px",
                borderRadius: "10px",
                fontSize: "12px",
              }}
            >
              {statusBuka.status}
            </span>
          </div>

          <div className="info-row">
            <span className="label">Fasilitas :</span> ✔ Wifi ✔ Live Music ✔
            Outdoor
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
              <div key={day} className="day-time-row">
                <span className="day-label">
                  {day.charAt(0).toUpperCase() + day.slice(1)}:
                </span>
                <input
                  type="time"
                  data-day={day}
                  data-field="buka"
                  value={formData.jamOperasional[day].buka}
                  onChange={handleInputChange}
                />
                <span> - </span>
                <input
                  type="time"
                  data-day={day}
                  data-field="tutup"
                  value={formData.jamOperasional[day].tutup}
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
