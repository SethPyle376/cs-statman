import Vue from 'vue'
import VueApollo from 'vue-apollo'
// import ApolloClient from 'apollo-boost'
import { createApolloClient } from 'vue-cli-plugin-apollo/graphql-client'

Vue.use(VueApollo)

const AUTH_TOKEN = 'DAYLY_TOKEN'

const httpEndpoint = process.env.VUE_APP_API_ENDPOINT || 'http://localhost:8082/graphql'

const defaultOptions = {
  // You can use `https` for secure connection (recommended in production)
  httpEndpoint,
  // You can use `wss` for secure connection (recommended in production)
  // Use `null` to disable subscriptions
  // wsEndpoint: process.env.VUE_APP_GRAPHQL_WS || 'ws://localhost:3000/graphql',
  // LocalStorage token
  tokenName: AUTH_TOKEN,
  // Enable Automatic Query persisting with Apollo Engine
  persisting: false,
  // Use websockets for everything (no HTTP)
  // You need to pass a `wsEndpoint` for this to work
  websocketsOnly: false,
  // Is being rendered on the server?
  ssr: false,

  // Override default apollo link
  // note: don't override httpLink here, specify httpLink options in the
  // httpLinkOptions property of defaultOptions.
  // link: myLink

  // Override default cache
  // cache: myCache

  // Override the way the Authorization header is set
  getAuth: () => {
    return 'Bearer: ' + localStorage.getItem(AUTH_TOKEN)
  }

  // Additional ApolloClient options
  // apollo: { ... }

  // Client local data (see apollo-link-state)
  // clientState: { resolvers: { ... }, defaults: { ... } }
}

export function createProvider (options = {}) {
  // const apolloClient = new ApolloClient({
  //   uri: 'http://localhost:8080/graphql'
  // })

  const { apolloClient } = createApolloClient({
    ...defaultOptions,
    ...options
  })

  const apolloProvider = new VueApollo({
    defaultClient: apolloClient,
    defaultOptions: {
      $query: {
        fetchPolicy: 'network-only'
      }
    }
  })

  return apolloProvider
}