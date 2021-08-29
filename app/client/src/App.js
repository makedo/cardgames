import React from 'react';
import Durak from "./games/Durak";
import { DndProvider } from 'react-dnd'
import { HTML5Backend } from 'react-dnd-html5-backend'

function App() {
  return <DndProvider backend={HTML5Backend}>
    <Durak />
  </DndProvider>
}

export default App;
