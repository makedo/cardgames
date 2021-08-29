import React from 'react';

import './style.css';

export default function Container({children}) {
  return <div className="card-container">
    {children}
  </div>
}
