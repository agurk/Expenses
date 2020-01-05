<template>
  <div id="app">
    <div class="header">
      <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="navbar-brand">Expenses Time Too!</div>
          <ul class="navbar-nav mr-auto">
              <li class="nav-item"><router-link to="/expenses" class="nav-link">Expenses</router-link></li>
              <li class="nav-item"><router-link to="/analysis" class="nav-link">Analysis</router-link></li>
              <li class="nav-item"><router-link to="/documents" class="nav-link">Documents</router-link></li>
          </ul>

        <form class="form-inline" action="/search">
          <input class="form-control mr-sm-2" name="description" type="search" placeholder="Search" aria-label="Search">
          <button class="btn" type="submit">Search</button>
        </form>
      </nav>
    </div>

      <router-view/>
  </div>
</template>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#app a { 
  color: #11557C;
}


</style>
