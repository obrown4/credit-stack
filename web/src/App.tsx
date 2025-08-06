import credit_icon from './assets/credit_icon.png'
import './App.css'

function App() {

  return (
    <>
      <div>
          <img src={credit_icon} className="logo" alt="Credit Stack" />
      </div>
      <h1>Credit Stack</h1>
      <div className="card"></div>
    </>
  )
}

export default App
