<template>
  <div id="app">
    <b-navbar toggleable="md">
      <b-navbar-brand>Expenses Time Too!</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
        <b-nav-item><router-link to="/expenses">Expenses</router-link></b-nav-item>
        <b-nav-item><router-link to="/analysis">Analysis</router-link></b-nav-item>
        <b-nav-item><router-link to="/assets">Assets</router-link></b-nav-item>
        <b-nav-item><router-link to="/documents">Documents</router-link></b-nav-item>
        <b-nav-item><router-link to="/search">Search</router-link></b-nav-item>
        <b-nav-item><router-link to="/config">Config</router-link></b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>

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
