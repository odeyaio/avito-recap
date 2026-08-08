import { setupServer } from "msw/node";

import { mockHandlers } from "../api/mocks/handlers";

export const server = setupServer(...mockHandlers);
