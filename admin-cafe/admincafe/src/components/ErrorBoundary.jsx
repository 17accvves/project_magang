import React from 'react';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, errorInfo) {
    console.error('Error Boundary caught:', error, errorInfo);
  }
  
  render() {
    if (this.state.hasError) {
      return (
        <div style={{ 
          padding: '40px', 
          border: '2px solid #ff6b6b', 
          background: '#ffe6e6',
          borderRadius: '12px',
          margin: '20px',
          textAlign: 'center',
          color: '#d63031'
        }}>
          <h3>⚠️ Terjadi Kesalahan</h3>
          <p>Component Ulasan sedang bermasalah: {this.state.error?.message}</p>
          <button 
            onClick={() => this.setState({ hasError: false, error: null })}
            style={{
              padding: '10px 20px',
              background: '#007bff',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              marginTop: '10px'
            }}
          >
            Coba Lagi
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}

export default ErrorBoundary;