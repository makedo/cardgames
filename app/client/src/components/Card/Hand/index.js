import React from 'react';
import './style.css';

export default function Hand({className, children}) {
  return <div className={["hand", className].filter(c => !!c).join(" ")}>
    {children}
  </div>
}
