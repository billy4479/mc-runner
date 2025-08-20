<script lang="ts">
  import { setMeOrError, getMeOrError } from "$lib/state.svelte";
  import LoginScreen from "./LoginScreen.svelte";
  import MainPage from "./MainPage.svelte";

  setMeOrError();
</script>

<main class="flex h-screen items-center justify-center gap-4">
  {#if getMeOrError() === null}
    Loading...
  {:else if getMeOrError().error}
    {#if getMeOrError().status == 401}
      <LoginScreen />
    {:else}
      <div>
        <span> An error has occurred </span>
        <pre>{JSON.stringify(getMeOrError(), null, 2)}</pre>
      </div>
    {/if}
  {:else}
    <MainPage />
  {/if}
</main>
