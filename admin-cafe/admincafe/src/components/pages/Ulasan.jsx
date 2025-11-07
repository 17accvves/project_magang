import React, { useEffect, useState } from "react";
import "./Ulasan.css";
import { FaStar, FaReply, FaPaperPlane, FaTimes, FaUser, FaImage, FaSync } from "react-icons/fa";

const API_URL = "http://localhost:8080/ulasan";
const ADMIN_API_URL = "http://localhost:8080/admin/ulasan";

const Ulasan = () => {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [replyingTo, setReplyingTo] = useState(null);
  const [replyText, setReplyText] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log("🔄 Fetching reviews from:", API_URL);
      const response = await fetch(API_URL);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      console.log("✅ Data received from backend:", data);
      setReviews(data);
      
    } catch (err) {
      console.error("❌ Gagal ambil data ulasan:", err);
      setError(`Tidak dapat terhubung ke server: ${err.message}`);
      setReviews([]);
    } finally {
      setLoading(false);
    }
  };

  const startReply = (reviewId) => {
    setReplyingTo(reviewId);
    setReplyText("");
  };

  const cancelReply = () => {
    setReplyingTo(null);
    setReplyText("");
  };

  const submitReply = async (reviewId) => {
    if (!replyText.trim()) {
      alert("Balasan tidak boleh kosong!");
      return;
    }

    setSaving(true);
    
    try {
      console.log("📤 Sending reply for review:", reviewId);
      
      const response = await fetch(`${ADMIN_API_URL}/${reviewId}/reply`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          reply: replyText
        }),
      });

      if (!response.ok) {
        throw new Error(`Gagal mengirim balasan: ${response.status}`);
      }

      const result = await response.json();
      console.log("✅ Reply successful:", result);
      
      // Refresh data setelah berhasil
      fetchReviews();
      setReplyingTo(null);
      setReplyText("");
      alert("Balasan berhasil dikirim!");
      
    } catch (error) {
      console.error("❌ Gagal mengirim balasan:", error);
      alert("Gagal mengirim balasan! Periksa koneksi backend.");
    } finally {
      setSaving(false);
    }
  };

  const deleteReply = async (reviewId) => {
    if (!window.confirm("Hapus balasan ini?")) return;

    try {
      console.log("🗑️ Deleting reply for review:", reviewId);
      
      const response = await fetch(`${ADMIN_API_URL}/${reviewId}/reply`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(`Gagal menghapus balasan: ${response.status}`);
      }

      const result = await response.json();
      console.log("✅ Delete reply successful:", result);
      
      // Refresh data setelah berhasil
      fetchReviews();
      alert("Balasan berhasil dihapus!");
      
    } catch (error) {
      console.error("❌ Gagal menghapus balasan:", error);
      alert("Gagal menghapus balasan! Periksa koneksi backend.");
    }
  };

  const renderStars = (rating) => {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(
        <FaStar 
          key={i} 
          color={i <= rating ? "#ffca2c" : "#e4e5e9"} 
        />
      );
    }
    return stars;
  };

  const testBackendConnection = async () => {
    try {
      setError("");
      const response = await fetch(API_URL);
      if (response.ok) {
        alert("✅ Backend terhubung dengan baik!");
        fetchReviews();
      } else {
        throw new Error(`Status: ${response.status}`);
      }
    } catch (err) {
      setError(`❌ Backend tidak terhubung: ${err.message}`);
    }
  };

  if (loading) {
    return (
      <div className="ulasan-container">
        <div className="ulasan-header">
          <h2 className="ulasan-title">Ulasan Pengunjung</h2>
          <button className="refresh-btn" onClick={fetchReviews}>
            <FaSync /> Refresh
          </button>
        </div>
        <div className="loading">
          <FaSync className="spinner" /> Memuat ulasan...
        </div>
      </div>
    );
  }

  return (
    <div className="ulasan-container">
      <div className="ulasan-header">
        <h2 className="ulasan-title">Ulasan Pengunjung</h2>
        <div className="header-actions">
          <button className="test-btn" onClick={testBackendConnection}>
            Test Koneksi
          </button>
          <button className="refresh-btn" onClick={fetchReviews}>
            <FaSync /> Refresh
          </button>
        </div>
      </div>
      
      {error && (
        <div className="error-warning">
          ⚠️ {error}
        </div>
      )}
      
      <div className="ulasan-stats">
        <div className="stat-card">
          <span className="stat-number">{reviews.length}</span>
          <span className="stat-label">Total Ulasan</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">
            {reviews.filter(r => r.rating >= 4).length}
          </span>
          <span className="stat-label">Ulasan Positif</span>
        </div>
        <div className="stat-card">
          <span className="stat-number">
            {reviews.filter(r => r.reply).length}
          </span>
          <span className="stat-label">Sudah Dibalas</span>
        </div>
      </div>
      
      <div className="ulasan-grid">
        {reviews.length === 0 ? (
          <div className="no-data">
            <FaUser size={48} color="#7a5d3a" />
            <p>Belum ada ulasan dari pengunjung</p>
            <p>Data akan muncul setelah ada pengunjung yang memberikan ulasan</p>
          </div>
        ) : (
          reviews.map((review) => (
            <div key={review.id} className="ulasan-card">
              {/* Header dengan Avatar */}
              <div className="ulasan-header-card">
                <div className="user-info">
                  <img 
                    src={review.avatar} 
                    alt={review.name}
                    className="user-avatar"
                    onError={(e) => {
                      e.target.src = `https://i.pravatar.cc/100?u=${review.id}`;
                    }}
                  />
                  <div className="user-details">
                    <h4 className="user-name">{review.name}</h4>
                    <div className="ulasan-meta">
                      <span className="ulasan-date">{review.date}</span>
                      <div className="ulasan-rating">
                        {renderStars(review.rating)}
                        <span className="rating-value">{review.rating}/5</span>
                      </div>
                    </div>
                  </div>
                </div>
                
                {/* Tombol Balas */}
                {!review.reply && (
                  <button 
                    className="balas-btn"
                    onClick={() => startReply(review.id)}
                    disabled={replyingTo === review.id}
                  >
                    <FaReply size={14} />
                    Balas
                  </button>
                )}
              </div>

              {/* Isi Ulasan */}
              <div className="ulasan-content">
                <p className="ulasan-text">"{review.text}"</p>
                
                {/* Gambar dari Ulasan */}
                {review.image && (
                  <div className="ulasan-image-container">
                    <img 
                      src={review.image} 
                      alt="Ulasan pengunjung" 
                      className="ulasan-image"
                      onError={(e) => {
                        e.target.style.display = 'none';
                      }}
                    />
                    <div className="image-badge">
                      <FaImage size={12} />
                      Foto dari pengunjung
                    </div>
                  </div>
                )}
              </div>

              {/* Balasan Admin */}
              {review.reply && (
                <div className="balasan-admin">
                  <div className="balasan-header">
                    <div className="admin-badge">
                      <strong>👨‍💼 Balasan Admin</strong>
                    </div>
                    <button 
                      className="delete-reply"
                      onClick={() => deleteReply(review.id)}
                      title="Hapus balasan"
                    >
                      <FaTimes size={12} />
                    </button>
                  </div>
                  <p className="balasan-teks">"{review.reply}"</p>
                </div>
              )}

              {/* Form Balasan */}
              {replyingTo === review.id && (
                <div className="reply-section">
                  <div className="reply-header">
                    <strong>Balas Ulasan dari {review.name}</strong>
                  </div>
                  <textarea
                    className="reply-input"
                    value={replyText}
                    onChange={(e) => setReplyText(e.target.value)}
                    placeholder="Tulis balasan Anda sebagai admin..."
                    rows="3"
                  />
                  <div className="reply-actions">
                    <button
                      className="save-reply-btn"
                      onClick={() => submitReply(review.id)}
                      disabled={saving || !replyText.trim()}
                    >
                      <FaPaperPlane size={12} />
                      {saving ? "Mengirim..." : "Kirim Balasan"}
                    </button>
                    <button
                      className="cancel-reply-btn"
                      onClick={cancelReply}
                      disabled={saving}
                    >
                      <FaTimes size={12} />
                      Batal
                    </button>
                  </div>
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {/* Debug Info */}
      <div className="debug-info">
        <details>
          <summary>ℹ️ Info Debug</summary>
          <p><strong>API URL:</strong> {API_URL}</p>
          <p><strong>Total Data:</strong> {reviews.length} ulasan</p>
          <p><strong>Backend Status:</strong> {error ? '❌ Offline' : '✅ Online'}</p>
          <button onClick={() => console.log('Reviews data:', reviews)}>
            Lihat Data di Console
          </button>
        </details>
      </div>
    </div>
  );
};

export default Ulasan;