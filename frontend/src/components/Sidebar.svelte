<script>
  import { createEventDispatcher } from 'svelte'

  export let uiPackage = null

  const dispatch = createEventDispatcher()

  function selectControl(control) {
    dispatch('selectControl', control)
  }

  function renderControls(controls, depth = 0) {
    if (!controls) return []

    return controls.map((control, idx) => ({
      id: `${depth}-${idx}`,
      control,
      depth,
      hasChildren: control.children && control.children.length > 0,
    }))
  }

  let expandedControls = new Set()

  function toggleExpand(id) {
    if (expandedControls.has(id)) {
      expandedControls.delete(id)
    } else {
      expandedControls.add(id)
    }
    expandedControls = expandedControls
  }

  $: controls = uiPackage?.Prototypes || []
</script>

<div class="sidebar">
  <div class="sidebar-header">
    <h2>Controls</h2>
  </div>
  <div class="sidebar-content">
    {#if controls.length === 0}
      <div class="empty-state">
        No file loaded. Open a file to see controls.
      </div>
    {:else}
      <div class="control-tree">
        {#each controls as control (control.name)}
          <div class="control-item">
            <div
              class="control-label"
              on:click={() => selectControl(control)}
            >
              {#if control.children && control.children.length > 0}
                <span
                  class="expand-icon"
                  on:click={(e) => {
                    e.stopPropagation()
                    toggleExpand(control.name)
                  }}
                >
                  {expandedControls.has(control.name) ? '▼' : '▶'}
                </span>
              {:else}
                <span class="expand-icon empty">•</span>
              {/if}
              <span class="control-name">{control.name}</span>
              <span class="control-class">{control.class}</span>
            </div>
            {#if expandedControls.has(control.name) && control.children}
              <div class="children">
                {#each control.children as child (child.name)}
                  <div class="child-item">
                    <span class="child-name" on:click={() => selectControl(child)}>
                      • {child.name}
                    </span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .sidebar {
    display: flex;
    flex-direction: column;
    background: #252526;
    border-right: 1px solid #3e3e3e;
    overflow: hidden;
  }

  .sidebar-header {
    padding: 0.75rem;
    border-bottom: 1px solid #3e3e3e;
    background: #2d2d2d;
  }

  .sidebar-header h2 {
    margin: 0;
    font-size: 0.95rem;
    font-weight: 600;
    color: #cccccc;
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem 0;
  }

  .empty-state {
    padding: 1rem;
    color: #888;
    font-size: 0.9rem;
    text-align: center;
  }

  .control-tree {
    display: flex;
    flex-direction: column;
  }

  .control-item {
    user-select: none;
  }

  .control-label {
    display: flex;
    align-items: center;
    padding: 0.4rem 0.5rem;
    cursor: pointer;
    color: #d4d4d4;
    font-size: 0.9rem;
    transition: background 0.15s;
  }

  .control-label:hover {
    background: #3e3e3e;
  }

  .expand-icon {
    display: inline-block;
    width: 1.2rem;
    text-align: center;
    font-size: 0.8rem;
    color: #888;
    cursor: pointer;
  }

  .expand-icon.empty {
    opacity: 0.3;
  }

  .control-name {
    font-weight: 500;
    margin-right: 0.5rem;
  }

  .control-class {
    font-size: 0.8rem;
    color: #888;
    margin-left: auto;
  }

  .children {
    display: flex;
    flex-direction: column;
    margin-left: 1rem;
    border-left: 1px solid #3e3e3e;
    padding-left: 0.5rem;
  }

  .child-item {
    padding: 0.3rem 0;
  }

  .child-name {
    display: block;
    padding: 0.2rem 0.5rem;
    color: #a0a0a0;
    font-size: 0.85rem;
    cursor: pointer;
    border-radius: 2px;
    transition: background 0.15s;
  }

  .child-name:hover {
    background: #3e3e3e;
  }

  /* Scrollbar styling */
  .sidebar-content::-webkit-scrollbar {
    width: 8px;
  }

  .sidebar-content::-webkit-scrollbar-track {
    background: transparent;
  }

  .sidebar-content::-webkit-scrollbar-thumb {
    background: #464647;
    border-radius: 4px;
  }

  .sidebar-content::-webkit-scrollbar-thumb:hover {
    background: #5a5a5a;
  }
</style>
