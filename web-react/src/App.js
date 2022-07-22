import React from 'react';
import './App.css';
import Recipe from './Recipe';
class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      recipes: []
    }
    this.getRecipes();
  }

  render() {
    return (<div>
      {this.state.recipes.map((recipe, index) => (
        <Recipe recipe={recipe} />
      ))}
    </div>);
  }

  getRecipes() {
    fetch('http://localhost:8080/recipes')
      .then(response => response.json())
      .then(data => this.setState({ recipes: data }))
  }

}

export default App;
