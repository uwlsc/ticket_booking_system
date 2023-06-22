import 'package:flutter/material.dart';
import '../common/color_constants.dart';
import '../widgets/app_widget.dart';
import '../widgets/text_button_widget.dart';
import '../widgets/video_clipper.dart';
import '../widgets/video_clipper2.dart';
import 'package:video_player/video_player.dart';

/*
Title: BookingScreen
Purpose:BookingScreen
Created By:Kalpesh Khandla
*/

class BookingScreen extends StatefulWidget {
  final VideoPlayerController moviePlayerController;
  final VideoPlayerController reflectionPlayerController;
  final String movieName;

  BookingScreen({
    Key? key,
    required this.moviePlayerController,
    required this.reflectionPlayerController,
    required this.movieName,
  }) : super(key: key);

  @override
  _BookingScreenState createState() => _BookingScreenState();
}

class _BookingScreenState extends State<BookingScreen>
    with TickerProviderStateMixin {
  Size get _size => MediaQuery.of(context).size;
  List<AnimationController> _dateSelectorACList = [];
  List<Animation<double>> _dateSelectorTweenList = [];

  List<AnimationController> _timeSelectorACList = [];
  List<Animation<double>> _timeSelectorTweenList = [];

  late AnimationController _dateBackgroundAc;
  late Animation<double> _dateBackgroundTween;

  late AnimationController _cinemaScreenAc;
  late Animation<double> _cinemaScreenTween;

  late AnimationController _reflectionAc;
  late Animation<double> _reflectionTween;

  late AnimationController _payButtonAc;
  late Animation<double> _payButtonTween;

  late AnimationController _cinemaChairAc;
  late Animation<double> _cinemaChairTween;

  int _dateIndexSelected = 1;
  int _timeIndexSelected = 1;
  var _chairStatus = [
    [0, 3, 3, 1, 2, 2, 1, 1 ,1, 1, 3, 2, 1, 2, 2, 1, 1 ,1, 1, 0 ],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 2],
    [1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 2, 2, 2, 1, 2, 2, 2, 2, 1, 1, 2, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 1, 1, 1, 2, 2, 1, 1, 1, 0, 0, 0, 2, 2, 2, 2, 1, 1, 2, 2],
    [2, 1, 1, 1, 2, 2, 1, 1, 1, 0, 0, 0, 2, 2, 2, 2, 1, 1, 2, 2],
    [3, 3, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0 , 0, 1, 1, 0,1, 1,0]
  ];

  @override
  void initState() {
    super.initState();
    widget.moviePlayerController.setLooping(true);
    widget.reflectionPlayerController.setLooping(true);
    widget.moviePlayerController.play();
    widget.reflectionPlayerController.play();

    initializeAnimation();
  }

  void initializeAnimation() {
    for (int i = 0; i < 7; i++) {
      _dateSelectorACList.add(AnimationController(
          vsync: this, duration: Duration(milliseconds: 500)));
      _dateSelectorTweenList.add(Tween<double>(begin: 1000, end: 0)
          .chain(CurveTween(curve: Curves.easeOutCubic))
          .animate(_dateSelectorACList[i]));
      Future.delayed(Duration(milliseconds: i * 50 + 170), () {
        _dateSelectorACList[i].forward();
      });
    }

    _dateBackgroundAc =
        AnimationController(vsync: this, duration: Duration(milliseconds: 700));
    _dateBackgroundTween = Tween<double>(begin: 1000, end: 0)
        .chain(CurveTween(curve: Curves.easeOutCubic))
        .animate(_dateBackgroundAc);
    Future.delayed(Duration(milliseconds: 150), () {
      _dateBackgroundAc.forward();
    });

    for (int i = 0; i < 3; i++) {
      _timeSelectorACList.add(AnimationController(
          vsync: this, duration: Duration(milliseconds: 500)));
      _timeSelectorTweenList.add(Tween<double>(begin: 1000, end: 0)
          .chain(CurveTween(curve: Curves.easeOutCubic))
          .animate(_timeSelectorACList[i]));
      Future.delayed(Duration(milliseconds: i * 25 + 100), () {
        _timeSelectorACList[i].forward();
      });
    }

    _cinemaScreenAc = AnimationController(
        vsync: this, duration: Duration(milliseconds: 2000));
    _cinemaScreenTween = Tween<double>(begin: 0, end: 1)
        .chain(CurveTween(curve: Curves.easeInOutQuart))
        .animate(_cinemaScreenAc);
    Future.delayed(Duration(milliseconds: 800), () {
      _cinemaScreenAc.forward();
    });

    // initialize reflection
    _reflectionAc =
        AnimationController(vsync: this, duration: Duration(milliseconds: 500));
    _reflectionTween = Tween<double>(begin: 0, end: 1)
        .chain(CurveTween(curve: Curves.easeInOutQuart))
        .animate(_reflectionAc);
    Future.delayed(Duration(milliseconds: 1800), () {
      _reflectionAc.forward();
    });

    // paybutton
    _payButtonAc = AnimationController(
        vsync: this, duration: Duration(milliseconds: 2000));
    _payButtonTween = Tween<double>(begin: -1, end: 0)
        .chain(CurveTween(curve: Curves.easeInOutQuart))
        .animate(_payButtonAc);
    Future.delayed(Duration(milliseconds: 800), () {
      _payButtonAc.forward();
    });

    // chair
    _cinemaChairAc = AnimationController(
        vsync: this, duration: Duration(milliseconds: 1600));
    _cinemaChairTween = Tween<double>(begin: -1, end: 0)
        .chain(CurveTween(curve: Curves.ease))
        .animate(_cinemaChairAc);
    Future.delayed(Duration(milliseconds: 1200), () {
      _cinemaChairAc.forward();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      body: Column(
        children: <Widget>[
          _appBar(),
          _dateSelector(),
          _timeSelector(),
          _cinemaRoom(),
          _payButton()
        ],
      ),
    );
  }

  Widget _cinemaRoom() {
    return Expanded(
        flex: 47,
        child: Stack(
          alignment: Alignment.topCenter,
          children: <Widget>[
            Container(
              width: _size.width,
            ),
            Container(
              padding: EdgeInsets.only(top: 18),
              child: Container(
                width: MediaQuery.of(context).size.width * .8,
                height: 800,
                child: Center(
                  child: widget.reflectionPlayerController.value.isInitialized
                      ? AnimatedBuilder(
                          animation: _reflectionAc,
                          builder: (ctx, child) {
                            return Opacity(
                              opacity: _reflectionTween.value,
                              child: child,
                            );
                          },
                          child: AspectRatio(
                            aspectRatio: .5,
                            child:
                                VideoPlayer(widget.reflectionPlayerController),
                          ),
                        )
                      : Container(),
                ),
              ),
            ),
            Positioned(
              top: 48,
              child: ClipPath(
                clipper: VideoClipper2(),
                child: AnimatedBuilder(
                  animation: _reflectionAc,
                  builder: (ctx, child) {
                    return Opacity(
                      opacity: _reflectionTween.value,
                      child: child,
                    );
                  },
                  child: Container(
                    height: 40,
                    width: MediaQuery.of(context).size.width * .8,
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [Colors.grey[300]!, Colors.transparent],
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        stops: [0, 1],
                      ),
                    ),
                  ),
                ),
              ),
            ),
            Positioned(
              bottom: _size.height * .02,
              child: AnimatedBuilder(
                animation: _cinemaChairTween,
                builder: (ctx, child) {
                  return Transform.translate(
                    offset: Offset(0, _cinemaChairTween.value * 100),
                    child: Opacity(
                        opacity: _cinemaChairTween.value + 1, child: child),
                  );
                },
                child: Container(
                  width: _size.width,
                  child: _chairList(),
                ),
              ),
            ),
            Positioned(
              top: 0,
              child: Container(
                width: MediaQuery.of(context).size.width * .8,
                child: Column(
                  children: <Widget>[
                    Center(
                      child: widget.moviePlayerController.value.isInitialized
                          ? AnimatedBuilder(
                              animation: _cinemaScreenAc,
                              builder: (ctx, child) {
                                double perspective =
                                    0.004 * _cinemaScreenTween.value;
                                // return ClipPath(
                                //     clipper: VideoClipper(
                                //         curveValue: monitorTween.value),
                                //     child: child);
                                return AspectRatio(
                                  aspectRatio: 16 / 9,
                                  child: Transform(
                                    alignment: Alignment.topCenter,
                                    transform: Matrix4.identity()
                                      ..setEntry(3, 2, perspective)
                                      ..rotateX(_cinemaScreenTween.value),
                                    child: ClipPath(
                                      clipper: VideoClipper(
                                        curveValue: _cinemaScreenTween.value,
                                      ),
                                      child: child,
                                    ),
                                  ),
                                );
                              },
                              child: VideoPlayer(widget.moviePlayerController),
                            )
                          : Container(),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ));
  }

  Widget _payButton() {
    return Expanded(
      flex: 13,
      child: AnimatedBuilder(
        animation: _payButtonAc,
        builder: (ctx, child) {
          double opacity() {
            if (_payButtonTween.value + 1 < 0.2) {
              return (_payButtonTween.value + 1) * 5;
            }
            return 1;
          }

          return Transform.translate(
            offset: Offset(0, _payButtonTween.value * 200),
            child: Opacity(opacity: opacity(), child: child),
          );
        },
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: <Widget>[
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 32),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: <Widget>[
                  _chairCategory(Colors.white, "FREE"),
                  _chairCategory(ColorConstants.kPrimaryColor, "YOURS"),
                  _chairCategory(Colors.grey[700]!, "RESERVED"),
                  _chairCategory(Colors.red[800]!, "NOT AVAILABLE"),
                ],
              ),
            ),
            Container(
              margin: EdgeInsets.only(left: 32, right: 32, bottom: 8),
              child: TextButtonWidget(
                btnTxt: "Pay",
                btnBackColor: ColorConstants.kPrimaryColor,
                textColor: ColorConstants.kBlackColor,
                btnOnTap: () {
                  print("Pay now");
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _chairCategory(Color color, String category) {
    return Row(
      children: <Widget>[
        Container(
          height: 10,
          width: 10,
          margin: EdgeInsets.only(right: 8),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(2),
            color: color,
          ),
        ),
        Text(
          category,
          style: Theme.of(context).textTheme.caption!.copyWith(
                fontSize: 12,
                color: ColorConstants.kGreyColor,
                fontFamily: "Bebas-Neue",
              ),
        ),
      ],
    );
  }

  Widget _chairList() {
    // 0 is null
    // 1 is free
    // 2 is reserved
    // 3 is notavailable
    // 4 is yours

    return Container(
      child: Column(
        children: <Widget>[
          for (int i = 0; i < 6; i++)
            Container(
              margin: EdgeInsets.only(top: i == 3 ? _size.height * .02 : 0),
              child: Row(
                children: <Widget>[
                  for (int x = 0; x < 9; x++)
                    Expanded(
                      flex: x == 0 || x == 8 ? 2 : 1,
                      child: x == 0 ||
                              x == 8 ||
                              (i == 0 && x == 1) ||
                              (i == 0 && x == 7) ||
                              (i == 3 && x == 1) ||
                              (i == 3 && x == 7) ||
                              (i == 5 && x == 1) ||
                              (i == 5 && x == 7)
                          ? Container()
                          : GestureDetector(
                              onTap: () {
                                if (_chairStatus[i][x - 1] == 1) {
                                  setState(() {
                                    _chairStatus[i][x - 1] = 4;
                                  });
                                } else if (_chairStatus[i][x - 1] == 4) {
                                  setState(() {
                                    _chairStatus[i][x - 1] = 1;
                                  });
                                }
                              },
                              child: Container(
                                height: _size.width / 11 - 10,
                                margin: EdgeInsets.all(5),
                                child: _chairStatus[i][x - 1] == 1
                                    ? AppWidget.whiteChair()
                                    : _chairStatus[i][x - 1] == 2
                                        ? AppWidget.greyChair()
                                        : _chairStatus[i][x - 1] == 3
                                            ? AppWidget.redChair()
                                            : AppWidget.yellowChair(),
                              ),
                            ),
                    ),
                ],
              ),
            ),
        ],
      ),
    );
  }

  Widget _timeSelector() {
    var time = [
      ["01.30", 1, 45],
      ["06.30", 2, 50],
      ["10.30", 3, 40]
    ];

    return Expanded(
      flex: 17,
      child: Container(
        margin: EdgeInsets.symmetric(vertical: _size.height * .035),
        child: ListView.builder(
          scrollDirection: Axis.horizontal,
          itemCount: 3,
          itemBuilder: (ctx, index) {
            return AnimatedBuilder(
              animation: _timeSelectorACList[index],
              builder: (ctx, child) {
                return Transform.translate(
                  offset: Offset(_timeSelectorTweenList[index].value, 0),
                  child: child,
                );
              },
              child: Container(
                margin: EdgeInsets.only(left: index == 0 ? 32 : 0),
                child: GestureDetector(
                  onTap: () {
                    setState(() {
                      _timeIndexSelected = index;
                    });
                  },
                  child: _timeItem(time[index][0] as String,time[index][1] as int, time[index][2] as int,
                      index == _timeIndexSelected ? true : false),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _timeItem(String time, int hall, int price, bool active) {
    return Container(
      margin: EdgeInsets.symmetric(horizontal: 12),
      padding: EdgeInsets.symmetric(horizontal: 16),
      decoration: BoxDecoration(
        border: Border.all(
          color: active
              ? ColorConstants.kPrimaryColor
              : ColorConstants.kWhiteColor,
          width: 1,
        ),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          RichText(
            text: TextSpan(
              text: time,
              style: Theme.of(context).textTheme.caption!.copyWith(
                    fontSize: 20,
                    fontWeight: FontWeight.w600,
                    color: active
                        ? ColorConstants.kPrimaryColor
                        : ColorConstants.kWhiteColor,
                  ),
              children: <TextSpan>[
                TextSpan(
                  text: ' PM',
                  style: Theme.of(context).textTheme.caption!.copyWith(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: active
                            ? ColorConstants.kPrimaryColor
                            : ColorConstants.kWhiteColor,
                      ),
                )
              ],
            ),
          ),
          Text(
            "Theatre${hall}( £${price})",
            style: Theme.of(context).textTheme.caption!.copyWith(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: ColorConstants.kGreyColor,
                ),
          )
        ],
      ),
    );
  }

  Widget _dateSelector() {
    DateTime currentDate = DateTime.now();

    return Expanded(
      flex: 13,
      child: Container(
        width: _size.width,
        padding: EdgeInsets.only(left: 32),
        child: Stack(
          alignment: Alignment.centerLeft,
          children: <Widget>[
            AnimatedBuilder(
              animation: _dateBackgroundAc,
              builder: (ctx, child) {
                return Transform.translate(
                  offset: Offset(_dateBackgroundTween.value, 0),
                  child: child,
                );
              },
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(.1),
                  borderRadius: BorderRadius.only(
                    topLeft: Radius.circular(12),
                    bottomLeft: Radius.circular(12),
                  ),
                ),
              ),
            ),
            Container(
              width: _size.width,
              child: ListView.builder(
                itemCount: 7,
                scrollDirection: Axis.horizontal,
                itemBuilder: (ctx, index) {
                  var date = currentDate.add(
                    Duration(days: index),
                  );

                  return AnimatedBuilder(
                    animation: _dateSelectorACList[index],
                    builder: (ctx, child) {
                      return Transform.translate(
                        offset: Offset(_dateSelectorTweenList[index].value, 0),
                        child: child,
                      );
                    },
                    child: GestureDetector(
                      onTap: () {
                        setState(() {
                          _dateIndexSelected = index;
                        });
                      },
                      child: Container(
                        padding: EdgeInsets.all(4),
                        margin: EdgeInsets.symmetric(
                          vertical: _size.height * .025,
                          horizontal: 12,
                        ),
                        width: 44,
                        decoration: BoxDecoration(
                          color: _dateIndexSelected == index
                              ? ColorConstants.kPrimaryColor
                              : Colors.transparent,
                          borderRadius: BorderRadius.circular(5),
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: <Widget>[
                            Text(
                              _dayFormat(date.weekday),
                              style: TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.w600,
                                color: index == _dateIndexSelected
                                    ? ColorConstants.kBlackColor
                                    : ColorConstants.kWhiteColor,
                              ),
                            ),
                            Text(
                              date.day.toString(),
                              style: TextStyle(
                                fontSize: 22,
                                fontWeight: FontWeight.w600,
                                height: 1,
                                color: index == _dateIndexSelected
                                    ? ColorConstants.kBlackColor
                                    : ColorConstants.kWhiteColor,
                              ),
                            )
                          ],
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _dayFormat(int dayWeek) {
    switch (dayWeek) {
      case 1:
        return "MO";
        break;
      case 2:
        return "TU";
        break;
      case 3:
        return "WE";
        break;
      case 4:
        return "TH";
        break;
      case 5:
        return "FR";
        break;
      case 6:
        return "SA";
        break;
      case 7:
        return "SU";

        break;
      default:
        return "MO";
    }
  }

  Widget _appBar() {
    return Expanded(
      flex: 8,
      child: Container(
        width: _size.width,
        padding: EdgeInsets.symmetric(vertical: 0),
        child: Stack(
          alignment: Alignment.center,
          children: <Widget>[
            Text(
              widget.movieName,
              style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold),
            ),
            Positioned(
              left: 24,
              child: Material(
                color: Colors.transparent,
                child: InkWell(
                  onTap: () {
                    Navigator.pop(context);
                  },
                  child: Icon(
                    Icons.arrow_back_ios,
                    color: Colors.white,
                    size: 18,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
