import { applyMiddleware, compose, createStore, Reducer, Store } from "redux";
import thunk from "redux-thunk";
import { combineReducers } from "redux";

export function noopReducer(
  state: Record<string, unknown> = {},
  _: Record<string, unknown>,
): Record<string, unknown> {
  return state;
}

export const rootReducer = combineReducers({
  base: noopReducer,
});

export function configureStore(
  initialState: { base: Record<string, unknown> },
  reducer: Reducer = rootReducer,
): Store<{ base: Record<string, unknown> }> {
  const middleware = [];
  if (typeof window !== "undefined") {
    middleware.push(thunk);
  }

  const enhancers = [applyMiddleware(...middleware)];

  if (
    typeof window !== "undefined" &&
    window &&
    // @ts-ignore
    window.__REDUX_DEVTOOLS_EXTENSION__
  ) {
    // @ts-ignore
    enhancers.push(window.__REDUX_DEVTOOLS_EXTENSION__());
  }

  return createStore(reducer, initialState, compose(...enhancers));
}
