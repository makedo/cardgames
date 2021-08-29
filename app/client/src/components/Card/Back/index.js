import React from 'react';

import './style.css';

export default function Back({className}) {
  return <div className={[className, "card", "back"].filter(c => !!c).join(' ')} />
}
