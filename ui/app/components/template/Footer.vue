<template>
  <footer class="pc-footer">
    <div class="footer-wrapper container-fluid">
      <div class="row">
        <div class="col-sm my-1">
          <p class="m-0">
            Copyright &copy; {{ currentYear }} All Right Reserved |
            <a href="https://www.kemenpppa.go.id" target="_blank" rel="noopener">
              Kemen PPPA
            </a>
          </p>
        </div>
        <div class="col-auto my-1">
          <ul class="list-inline footer-link mb-0">
            <li class="list-inline-item">
              <NuxtLink to="/">Arsip Record Center</NuxtLink>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>
