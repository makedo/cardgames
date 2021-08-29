import React from 'react';
import './style.css';

export const SPADES = 1;
export const CLUBS = 2;
export const HEARTS = 3;
export const DIAMONDS = 4;

export const Clubs    = () => <span className="black">&clubs;</span>;
export const Hearts   = () => <span className="red">&hearts;</span>;
export const Diamonds = () => <span className="red">&diams;</span>;
export const Spades   = () => <span className="black">&spades;</span>;

export default function Suite({ suite }) {
  switch (suite) {
    case SPADES:
      return <Spades />
    case CLUBS:
      return <Clubs />
    case HEARTS:
      return <Hearts />
    case DIAMONDS:
      return <Diamonds />
    default:
      throw new Error('Invalid suite');
  }
}
