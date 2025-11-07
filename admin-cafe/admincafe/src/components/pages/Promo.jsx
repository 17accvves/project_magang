import React, { useState, useEffect, useMemo } from "react";
import "./Promo.css";
import { FaChevronUp, FaChevronDown, FaSpinner, FaTags, FaPercent, FaCalendarAlt, FaShoppingBag } from "react-icons/fa";

const Promo = () => {
  const [promos, setPromos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showScrollTop, setShowScrollTop] = useState(true);

  useEffect(() => {
    fetchPromoMenus();
  }, []);

  const fetchPromoMenus = async () => {
    try {
      setLoading(true);
      console.log('🔄 Fetching promo menus...');
      
      const response = await fetch('http://localhost:8080/api/v1/promos');
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const result = await response.json();
      console.log('📦 API Response:', result);
      
      if (result.error) {
        throw new Error(result.message);
      }
      
      const promoData = result.data || [];
      console.log('✅ Promo data received:', promoData.length, 'items');
      
      setPromos(promoData);
      
    } catch (err) {
      console.error('❌ Error fetching promo menus:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const totalRevenue = useMemo(() => {
    return promos.reduce((sum, promo) => {
      const simulatedRevenue = Math.floor((promo.price * promo.discount) / 100);
      return sum + (simulatedRevenue > 0 ? Math.floor(simulatedRevenue / 10000) : 0);
    }, 0);
  }, [promos]);

  const toggleScroll = () => {
    setShowScrollTop(!showScrollTop);
    window.scrollTo({
      top: showScrollTop ? 0 : document.body.scrollHeight,
      behavior: "smooth",
    });
  };

  const formatDate = (dateString) => {
    if (!dateString) return "Selamanya";
    
    try {
      const date = new Date(dateString + 'T00:00:00');
      const options = { day: 'numeric', month: 'short', year: 'numeric' };
      return date.toLocaleDateString('id-ID', options);
    } catch (error) {
      return dateString;
    }
  };

  const handleImageError = (e) => {
    e.target.style.display = 'none';
  };

  if (loading) {
    return (
      <div className="promo-container">
        <div className="loading-state">
          <FaSpinner className="spinner" />
          <p>Memuat promo...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="promo-container">
        <div className="error-state">
          <h3>Gagal Memuat Promo</h3>
          <p>{error}</p>
          <button onClick={fetchPromoMenus} className="retry-btn">
            Coba Lagi
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="promo-container">
      {/* Sidebar Stats */}
      <div className="promo-sidebar">
        <div className="sidebar-header">
          <FaTags className="sidebar-icon" />
          <h3>Statistik Promo</h3>
        </div>
        
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-number">{promos.length}</div>
            <div className="stat-label">Promo Aktif</div>
          </div>
          <div className="stat-card">
            <div className="stat-number">{totalRevenue}</div>
            <div className="stat-label">Revenue (JT)</div>
          </div>
        </div>

        <div className="promo-list-section">
          <h4>Daftar Promo</h4>
          <div className="promo-list">
            {promos.map((promo) => (
              <div key={promo.id} className="promo-list-item">
                <FaPercent className="promo-list-icon" />
                <span className="promo-list-name">{promo.name}</span>
                <span className="promo-list-discount">{promo.discount}%</span>
              </div>
            ))}
            {promos.length === 0 && (
              <div className="empty-state">
                Tidak ada promo aktif
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="promo-content">
        <div className="content-header">
          <h1>Promo Spesial</h1>
          <p>Nikmati penawaran terbaik dari menu favorit Anda</p>
        </div>

        <div className="promo-grid">
          {promos.length === 0 ? (
            <div className="no-promos">
              <FaTags size={64} color="#d4a356" />
              <h3>Belum Ada Promo</h3>
              <p>Tambahkan diskon pada menu untuk menampilkan promo di sini</p>
            </div>
          ) : (
            promos.map((promo) => (
              <div key={promo.id} className="promo-card">
                <div className="card-header">
                  <div className="promo-badge">
                    <FaPercent className="badge-icon" />
                    {promo.discount}% OFF
                  </div>
                  <div className="status-badge active">
                    AKTIF
                  </div>
                </div>

                <div className="card-body">
                  {promo.img && (
                    <div className="image-container">
                      <img 
                        src={promo.img.startsWith('http') ? promo.img : `http://localhost:8080${promo.img}`} 
                        alt={promo.name}
                        onError={handleImageError}
                      />
                    </div>
                  )}
                  
                  <h3 className="menu-name">{promo.name}</h3>
                  
                  <div className="price-section">
                    <div className="price-row">
                      <span className="price-label">Normal</span>
                      <span className="old-price">Rp {promo.price?.toLocaleString('id-ID')}</span>
                    </div>
                    <div className="price-row">
                      <span className="price-label">Promo</span>
                      <span className="new-price">Rp {promo.discountedPrice?.toLocaleString('id-ID')}</span>
                    </div>
                    <div className="savings">
                      Hemat Rp {(promo.price - promo.discountedPrice)?.toLocaleString('id-ID')}
                    </div>
                  </div>

                  <div className="promo-details">
                    <div className="detail-item">
                      <FaCalendarAlt className="detail-icon" />
                      <div className="detail-content">
                        <span className="detail-label">Periode</span>
                        <span className="detail-value">
                          {formatDate(promo.startDate)} - {formatDate(promo.endDate)}
                        </span>
                      </div>
                    </div>
                    
                    <div className="detail-item">
                      <FaShoppingBag className="detail-icon" />
                      <div className="detail-content">
                        <span className="detail-label">Kategori</span>
                        <span className="detail-value">{promo.category}</span>
                      </div>
                    </div>
                  </div>

                  <div className="revenue-section">
                    <span>Revenue Estimate:</span>
                    <span className="revenue-value">
                      {Math.floor((promo.price * promo.discount) / 100000)}JT
                    </span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Scroll Button */}
        {promos.length > 2 && (
          <div className="scroll-icon" onClick={toggleScroll}>
            {showScrollTop ? <FaChevronUp /> : <FaChevronDown />}
          </div>
        )}
      </div>
    </div>
  );
};

export default Promo;