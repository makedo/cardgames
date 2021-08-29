import React, {useState, useEffect} from 'react';

import Container from "../../components/Card/Container";
import Face from "../../components/Card/Face";

import './style.css';

export default function Patience() {
  const [deck, setDeck] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch('/api/deck', {
        method: 'GET', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        const error = await response.text();
        console.error(error);
      } else {
        const data =  await response.json();
        setDeck(data);
      }
    }

    fetchData();
  }, []);

  return (deck ? <div className="patience">{deck.map((card) => <Container key={`${card.suite}${card.rank}`}><Face card={card} /></Container>)}</div> : 'Loading...');
}
