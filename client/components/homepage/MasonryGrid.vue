<script setup lang="ts">
import MasonryCard from "./MasonryCard.vue";
import { faker } from "@faker-js/faker";

function getRandomNumber(min: number, max: number) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function getNewImages() {
  return new Array(20).fill(null).map((_, i) => {
    const x = getRandomNumber(300, 500);
    const y = getRandomNumber(200, 500);
    const url = faker.image.urlLoremFlickr({
      category: "food",
      height: y,
      width: x,
    });

    return {
      url,
      username: faker.internet.displayName(),
      title: faker.lorem.words({ min: 1, max: 4 }),
      recipeId: i === 0 ? 1 : i + 2,
    };
  });
}

const initialImages = getNewImages();
initialImages[0].username = "testuser";

const images = ref(initialImages);
</script>

<template>
  <div class="px-16 min-h-screen" ref="container">
    <masonry-wall
      :items="images"
      :ssr-columns="1"
      :min-columns="1"
      :max-columns="3"
      :column-width="300"
      :gap="32"
    >
      <template #default="{ item, index }">
        <MasonryCard
          :image="item.url"
          :name="item.username"
          :title="item.title"
          :recipe-id="item.recipeId"
        />
      </template>
    </masonry-wall>
  </div>
  <div ref="end">end</div>
</template>

<style scoped></style>
