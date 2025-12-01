import React from 'react';
import './CafePage.css';

const CafeCard = ({ name, rank, reviews, logo, rating }) => (
  <div className="cafe-card">
    <div className="cafe-logo-area">
      <img src={logo} alt={name} className="cafe-logo-img" />
      <button className="img-nav left">‹</button>
      <button className="img-nav right">›</button>
    </div>
    <div className="cafe-info">
      <div className="cafe-header-row">
        <h3>{name}</h3>
        <button className="btn-icon-edit">✏️</button>
      </div>
      <p className="cafe-rank">#{rank} Bulan September</p>
      <div className="cafe-rating">
        <span className="stars">{'★'.repeat(Math.floor(rating))}</span>
        <span className="review-count">({reviews})</span>
      </div>
      <button className="btn-see-cafe">Lihat Cafe</button>
    </div>
  </div>
);

const CafePage = () => {
  const cafes = [
    { name: 'YOTTA.id', rank: 1, rating: 5, reviews: 131, logo: 'https://via.placeholder.com/150/2ecc71/ffffff?text=YOTTA' },
    { name: 'Fore Coffee', rank: 2, rating: 4, reviews: 120, logo: 'https://via.placeholder.com/150/ffffff/2ecc71?text=fore' },
    { name: 'Starbucks', rank: 3, rating: 4.5, reviews: 117, logo: 'https://via.placeholder.com/150/006241/ffffff?text=Starbucks' },
  ];

  return (
    <div className="cafe-content-padding cafe-theme"> 
      <div className="page-header-row">
        <h1 className="page-title">Cafe</h1>
        <button className="btn-add-cafe">
            <span className="plus">+</span> Daftarkan Cafe
        </button>
      </div>

      <div className="promo-banner">
        <div className="banner-nav left">‹</div>
        <div className="banner-content">
            <span className="banner-date">September 25-30</span>
            <h2 className="banner-title">Nikmati keseruan LIVE MUSIC <br /> dengan artis di daerah kamu!</h2>
            <p className="banner-subtitle">• CARI Spot Super ADMIN</p>
            <button className="btn-promo-action">Atur Promo</button>
        </div>
        <div className="banner-nav right">›</div>
      </div>

      <div className="cafe-grid">
        {cafes.map((cafe, index) => (
            <CafeCard key={index} {...cafe} />
        ))}
      </div>
    </div>
  );
};

export default CafePage;