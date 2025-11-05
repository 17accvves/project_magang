import React, { useEffect, useState } from "react";
import "./Ulasan.css";
import { FaStar } from "react-icons/fa";

const API_URL = "http://localhost:8080/ulasan";

const Ulasan = () => {
  const [reviews, setReviews] = useState([]);

  useEffect(() => {
    fetch(API_URL)
      .then((res) => res.json())
      .then((data) => setReviews(data))
      .catch((err) => console.error("Gagal ambil data ulasan:", err));
  }, []);

  const renderStars = (rating) => {
    const stars = [];
    const full = Math.floor(rating);
    const half = rating % 1 !== 0;
    for (let i = 0; i < full; i++) stars.push(<FaStar key={i} color="#ffca2c" />);
    if (half) stars.push(<FaStar key="half" color="#ffca2c" style={{ opacity: 0.5 }} />);
    return stars;
  };

  return (
    <div className="ulasan-container">
      <h2 className="ulasan-title">Ulasan Pengunjung</h2>
      <div className="ulasan-grid">
        {reviews.length === 0 ? (
          <p className="no-data">Belum ada ulasan.</p>
        ) : (
          reviews.map((r) => (
            <div key={r.id} className="ulasan-card">
              <div className="ulasan-date">{r.date}</div>
              <div className="ulasan-stars">{renderStars(r.rating)}</div>
              <p className="ulasan-text">“{r.text}”</p>
              {r.image && <img src={r.image} alt="ulasan" className="ulasan-image" />}
              {r.reply && (
                <div className="balasan-admin">
                  <strong>Balasan Admin :</strong> “{r.reply}”
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Ulasan;
