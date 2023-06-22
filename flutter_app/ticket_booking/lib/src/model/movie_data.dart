import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

Future<MovieModel> fetchMovie(String imdbId) async {
  final response = await http.get(Uri.parse('http://www.omdbapi.com/?i=$imdbId&apikey=337f9b3f&plot=full'));

  if (response.statusCode == 200) {
    // If the server returned a 200 OK response, parse the JSON.
    return MovieModel.fromJson(jsonDecode(response.body));
  } else {
    // If the server did not return a 200 OK response, throw an exception.
    throw Exception('Failed to load movie.');
  }
}

class MovieModel {
  String id;
  String name;
  List<String> genre;
  double rating;
  String storyLine;
  Image image;
  Image imageText;
  String videoClipPath;
  String videoClipReflectionPath;
  List<MovieCastModel> castList = [];

  MovieModel({
    required this.id,
    required this.genre,
    required this.name,
    required this.rating,
    required this.storyLine,
    required this.image,
    required this.imageText,
    required this.castList,
    required this.videoClipPath,
    required this.videoClipReflectionPath,
  });

  factory MovieModel.fromJson(Map<String, dynamic> json) {
    return MovieModel(
      id: json['imdbID'],
      genre: List<String>.from(json['Genre'].split(', ')),
      name: json['Title'],
      rating: double.parse(json['imdbRating']),
      storyLine: json['Plot'],
      image: Image.network(json['Poster']), // Create Image from URL
      imageText: Image.network('https://picsum.photos/200/50'),
      castList:  [
        MovieCastModel(
          name: "Robert Downey Jr",
          photo: Image.asset("assets/images/cast/robert.jpg"),
        ),
        MovieCastModel(
          name: "Tom Holland",
          photo: Image.asset("assets/images/cast/tom.jpg"),
        ),
        MovieCastModel(
          name: "Rami Malek",
          photo: Image.asset("assets/images/cast/rami.jpg"),
        ),
        MovieCastModel(
          name: "Selena Gomez",
          photo: Image.asset("assets/images/cast/selena.jpg"),
        )],
      videoClipPath: "assets/videos/mulanclip.mp4",
      videoClipReflectionPath: "assets/videos/dolittleclip-reflection.mp4",
    );
  }
}

class MovieCastModel {
  String name;
  Image photo;

  MovieCastModel({
    required this.name,
    required this.photo,
  });
}

class MovieData {
  List<MovieModel> movieList = [];

  MovieData() {
    movieList = [];
    loadDatabaseMovie();
    fetchMovies();
  }

  void loadDatabaseMovie() {
    movieList.add(
      MovieModel(
          id: "0",
          name: 'Dolittle',
          rating: 5.6,
          genre: ["Family", "Adventure"],
          storyLine:
              "Dr. John Dolittle lives in solitude behind the high walls of his lush manor in 19th-century England. His only companionship comes from an array of exotic animals that he speaks to on a daily basis. But when young Queen Victoria becomes gravely ill, the eccentric doctor and his furry friends embark on an epic adventure to a mythical island to find the cure.",
          image: Image.asset("assets/images/dolittle.jpg"),
          imageText: Image.asset("assets/images/dolittle-text.png"),
          videoClipPath: "assets/videos/dolittleclip.mp4",
          videoClipReflectionPath: "assets/videos/dolittleclip-reflection.mp4",
          castList: [
            MovieCastModel(
              name: "Robert Downey Jr",
              photo: Image.asset("assets/images/cast/robert.jpg"),
            ),
            MovieCastModel(
              name: "Tom Holland",
              photo: Image.asset("assets/images/cast/tom.jpg"),
            ),
            MovieCastModel(
              name: "Rami Malek",
              photo: Image.asset("assets/images/cast/rami.jpg"),
            ),
            MovieCastModel(
              name: "Selena Gomez",
              photo: Image.asset("assets/images/cast/selena.jpg"),
            )
          ]),
    );
    movieList.add(
      MovieModel(
        id: "1",
        name: 'Mulan',
        rating: 4.0,
        genre: ["Action", "Adventure"],
        storyLine:
            "Fearful that her ailing father will be drafted into the Chinese military, Mulan (Ming-Na Wen) takes his spot -- though, as a girl living under a patriarchal regime, she is technically unqualified to serve. She cleverly impersonates a man and goes off to train with fellow recruits. Accompanied by her dragon, Mushu (Eddie Murphy), she uses her smarts to help ward off a Hun invasion, falling in love with a dashing captain along the way.",
        image: Image.asset("assets/images/mulan.jpg"),
        imageText: Image.asset("assets/images/mulan-text.png"),
        videoClipPath: "assets/videos/mulanclip.mp4",
        videoClipReflectionPath: "assets/videos/mulanclip-reflection.mp4",
        castList: [
          MovieCastModel(
            name: "Liu Yifei",
            photo: Image.asset(
              "assets/images/cast/liu.jpg",
            ),
          ),
          MovieCastModel(
            name: "Donnie Yen",
            photo: Image.asset(
              "assets/images/cast/donnie.jpg",
            ),
          ),
          MovieCastModel(
            name: "Jet Li",
            photo: Image.asset(
              "assets/images/cast/jetli.jpg",
            ),
          ),
          MovieCastModel(
            name: "Gong Li",
            photo: Image.asset(
              "assets/images/cast/gongli.jpg",
            ),
          )
        ],
      ),
    );
    movieList.add(
      MovieModel(
        id: "2",
        name: 'Black Widow',
        rating: 7.0,
        genre: ["Action", "Adventure"],
        storyLine:
            "At birth the Black Widow (aka Natasha Romanova) is given to the KGB, which grooms her to become its ultimate operative. When the U.S.S.R. breaks up, the government tries to kill her as the action moves to present-day New York, where she is a freelance operative.",
        image: Image.asset(
          "assets/images/blackwidow.jpg",
        ),
        imageText: Image.asset(
          "assets/images/blackwidow-text.png",
        ),
        videoClipPath: "assets/videos/blackwidowclip.mp4",
        videoClipReflectionPath: "assets/videos/blackwidowclip-reflection.mp4",
        castList: [
          MovieCastModel(
            name: "Scarlett Johansson",
            photo: Image.asset(
              "assets/images/cast/scarlett.jpg",
            ),
          ),
          MovieCastModel(
            name: "Florance Pugh",
            photo: Image.asset(
              "assets/images/cast/florance.jpg",
            ),
          ),
          MovieCastModel(
            name: "David Harbour",
            photo: Image.asset(
              "assets/images/cast/david.jpg",
            ),
          ),
          MovieCastModel(
            name: "Rachel Weisz",
            photo: Image.asset(
              "assets/images/cast/rachel.jpg",
            ),
          )
        ],
      ),
    );
  }

  void fetchMovies() async {
    try {
      final movie1 = await fetchMovie('tt0816692');
      movieList.add(movie1);
    } catch (e) {
      print('Failed to fetch movie 1: $e');
    }

    try {
      final movie2 = await fetchMovie('tt4154796');
      movieList.add(movie2);
    } catch (e) {
      print('Failed to fetch movie 2: $e');
    }
  }
}
