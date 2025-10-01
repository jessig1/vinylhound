<script>
  export let message = "";
  export let kind = "info";
</script>

{#if message}
  <div class={"flash " + kind} role="status" aria-live="polite">
    {message}
  </div>
{/if}

<style>
  .flash {
    padding: 1rem 1.2rem;
    border-radius: 0.75rem;
    margin-bottom: 1.5rem;
    font-weight: 600;
  }

  .flash.success {
    background: #ecfdf5;
    color: #047857;
    border: 1px solid #a7f3d0;
  }

  .flash.error {
    background: #fef2f2;
    color: #b91c1c;
    border: 1px solid #fecaca;
  }

  .flash.info {
    background: #eff6ff;
    color: #1d4ed8;
    border: 1px solid #bfdbfe;
  }
</style>
