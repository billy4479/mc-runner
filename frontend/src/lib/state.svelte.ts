import { getMe } from "./api";

let meOrError = $state(null);

export async function setMeOrError() {
  meOrError = await getMe();
}

export function getMeOrError() {
  return meOrError;
}
