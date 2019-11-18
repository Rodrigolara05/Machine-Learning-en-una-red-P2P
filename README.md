# Sistema de recomendacion con Machine Learning en una red P2P

## Comenzando 🚀

El método de los k vecinos más cercanos (en inglés, k-nearest neighbors) es un método de clasificación supervisada 
(Aprendizaje, estimación basada en un conjunto de entrenamiento y prototipos) 
que sirve para estimar la función de densidad de las predictorasx por cada clase .

El sistema es un *clasificador de peliculas* para el usuario segun sus gustos anteriores (data historica), por lo que el sistema tendrá dos modos:
Entrenamiento, Prueba. Los datos solicitados de entrada son: Nombre, Clasificacion, Genero, Año y Gusto (Este ultimo solo si esta en modo entrenamiento)

Asimismo, el valor de K tendra que ser impar, como 3,5,7,9, etc y menor o igual a la cantidad de nodos entrenados.

## Arquitectura 📋

La arquitectura del sistema es una arquitecutra P2P,es decir, un nodo puede ser cliente y servidor a la vez. 

## Tecnicas 📌

* **Algoritmo de Búsqueda KNN* ** -  [A Asterisk Algorithm](https://es.wikipedia.org/wiki/K_vecinos_m%C3%A1s_pr%C3%B3ximos)
* **Peer To Peer* ** -  [P2P](https://es.wikipedia.org/wiki/Peer-to-peer)

## Versionado 📌

Usamos [Git](https://git-scm.com/) para el versionado.

## Autores ✒️

- Rodrigo Max Lara Camarena

* **Rodrigo Max Lara Camarena** -  [Rodrigo Lara](https://www.linkedin.com/in/rodrigolara05)

### Pre-requisitos 📋

Para poder trabajar con el siguiente proyecto debe de tener conocimientos de desarrollo de algoritmos sobre IA y la construccion de una arquitectura P2P.
Asi como tener conocimientos en el lenguaje de programación Go para la construcción de estos.
