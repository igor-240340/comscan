import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";

function MyButton() {
  return (
    <button>
      Test
    </button>
  );
}

function App() {
    return (
        <div id="App">
            <MyButton/>
        </div>
    )
}

export default App
