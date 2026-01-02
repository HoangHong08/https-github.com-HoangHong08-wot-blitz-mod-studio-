<script>
  import { onMount } from 'svelte'

  export let uiPackage = null
  export let selectedControl = null

  let previewCanvas

  onMount(() => {
    if (previewCanvas) {
      const ctx = previewCanvas.getContext('2d')
      if (ctx) {
        ctx.fillStyle = '#2d2d2d'
        ctx.fillRect(0, 0, previewCanvas.width, previewCanvas.height)
        
        // Draw grid pattern
        ctx.strokeStyle = '#3e3e3e'
        ctx.lineWidth = 1
        for (let i = 0; i < previewCanvas.width; i += 50) {
          ctx.beginPath()
          ctx.moveTo(i, 0)
          ctx.lineTo(i, previewCanvas.height)
          ctx.stroke()
        }
        for (let i = 0; i < previewCanvas.height; i += 50) {
          ctx.beginPath()
          ctx.moveTo(0, i)
          ctx.lineTo(previewCanvas.width, i)
          ctx.stroke()
        }
      }
    }
  })

  function drawControl(ctx, control, offsetX = 0, offsetY = 0) {
    if (!control || !control.size) return

    const x = (control.position?.x || 0) + offsetX
    const y = (control.position?.y || 0) + offsetY
    const width = control.size.x
    const height = control.size.y

    // Draw control background
    ctx.fillStyle = selectedControl?.name === control.name ? '#4CAF50' : '#3e3e3e'
    ctx.fillRect(x, y, width, height)

    // Draw border
    ctx.strokeStyle = selectedControl?.name === control.name ? '#66BB6A' : '#555'
    ctx.lineWidth = 1
    ctx.strokeRect(x, y, width, height)

    // Draw name
    if (width > 20 && height > 20) {
      ctx.fillStyle = '#d4d4d4'
      ctx.font = '12px monospace'
      ctx.fillText(control.name, x + 4, y + 16)
    }
  }

  function updatePreview() {
    if (!previewCanvas || !uiPackage) return

    const ctx = previewCanvas.getContext('2d')
    if (!ctx) return

    // Clear canvas
    ctx.fillStyle = '#1e1e1e'
    ctx.fillRect(0, 0, previewCanvas.width, previewCanvas.height)

    // Draw grid
    ctx.strokeStyle = '#3e3e3e'
    ctx.lineWidth = 1
    for (let i = 0; i < previewCanvas.width; i += 50) {
      ctx.beginPath()
      ctx.moveTo(i, 0)
      ctx.lineTo(i, previewCanvas.height)
      ctx.stroke()
    }
    for (let i = 0; i < previewCanvas.height; i += 50) {
      ctx.beginPath()
      ctx.moveTo(0, i)
      ctx.lineTo(previewCanvas.width, i)
      ctx.stroke()
    }

    // Draw all controls
    if (uiPackage.Prototypes) {
      uiPackage.Prototypes.forEach((proto) => {
        drawControl(ctx, proto)
        if (proto.children) {
          proto.children.forEach((child) => {
            drawControl(ctx, child, proto.position?.x || 0, proto.position?.y || 0)
          })
        }
      })
    }

    // Draw info
    ctx.fillStyle = '#888'
    ctx.font = '12px monospace'
    ctx.fillText('Preview - drag canvas to pan', 5, previewCanvas.height - 5)
  }

  $: if (previewCanvas) {
    updatePreview()
  }
</script>

<div class="preview">
  <div class="preview-header">
    <span>Preview</span>
  </div>
  <div class="preview-content">
    <canvas
      bind:this={previewCanvas}
      width={500}
      height={600}
      style="width: 100%; height: 100%; display: block;"
    />
  </div>
</div>

<style>
  .preview {
    display: flex;
    flex-direction: column;
    background: #1e1e1e;
    border-left: 1px solid #3e3e3e;
    overflow: hidden;
  }

  .preview-header {
    padding: 0.5rem 1rem;
    background: #2d2d2d;
    border-bottom: 1px solid #3e3e3e;
    font-size: 0.9rem;
    color: #cccccc;
    font-weight: 500;
  }

  .preview-content {
    flex: 1;
    overflow: auto;
    background: #1e1e1e;
  }

  .preview-content :global(canvas) {
    image-rendering: pixelated;
    background: #1e1e1e;
  }
</style>
