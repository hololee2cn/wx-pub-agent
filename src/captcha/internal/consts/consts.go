package consts

import "github.com/mojocn/base64Captcha"

const (
	Module = "shared-captcha"

	KeyTraceID = "x-trace-id"
	// 验证码redis key: Module + "_" + uuid
)

// 验证码类型
const (
	CaptchaTypeAudio   = "audio"
	CaptchaTypeString  = "string"
	CaptchaTypeMath    = "math"
	CaptchaTypeChinese = "chinese"
	CaptchaTypeDigit   = "digit"
)

const (
	// 干扰线设置, 分别为空心线, 细线, 正弦曲线
	OptionShowHollowLine = base64Captcha.OptionShowHollowLine
	OptionShowSlimeLine  = base64Captcha.OptionShowSlimeLine
	OptionShowSineLine   = base64Captcha.OptionShowSineLine
	OptionsShowAllLines  = OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine
)

const (
	// 音频类型的验证码默认语言
	// Language possible values for lang are "en", "ja", "ru", "zh".
	AudioDefaultLanguage = "zh"

	// 验证码图片默认 宽度, 高度, 字符数
	CaptchaDefaultWidth  = 150
	CaptchaDefaultHeight = 50
	CaptchaDefaultLength = 5

	// 噪音量
	CaptchaDefaultNoiseCount = 50

	// 默认有效期限180秒
	CaptchaDefaultMaxAge = 180

	// 验证码图片默认背景色: 全透明黑色
	CaptchaDefaultBgColor = "#00000000"
)

const (
	TxtNumbers            = "012346789"
	TxtAlphabet           = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
	TxtNumbersAndAlphabet = TxtNumbers + TxtAlphabet
	CommonlyUsedChinese   = "" +
		"的一了是我不在人们有" +
		"来他这上着个地到大里" +
		"说就去子得也和那要下" +
		"看天时过出小么起你都" +
		"把好还多没为又可家学" +
		"只以主会样年想生同老" +
		"中十从自面前头道它后" +
		"然走很像见两用她国动" +
		"进成回什边作对开而己" +
		"些现山民候经发工向事" +
		"命给长水几义三声于高" +
		"手知理眼志点心战二问" +
		"但身方实吃做叫当住听" +
		"革打呢真全才四已所敌" +
		"之最光产情路分总条白" +
		"话东席次亲如被花口放" +
		"儿常气五第使写军吧文" +
		"运再果怎定许快明行因" +
		"别飞外树物活部门无往" +
		"船望新带队先力完却站" +
		"代员机更九您每风级跟" +
		"笑啊孩万少直意夜比阶" +
		"连车重便斗马哪化太指" +
		"变社似士者干石满日决" +
		"百原拿群究各六本思解" +
		"立河村八难早论吗根共" +
		"让相研今其书坐接应关" +
		"信觉步反处记将千找争" +
		"领或师结块跑谁草越字" +
		"加脚紧爱等习阵怕月青" +
		"半火法题建赶位唱海七" +
		"女任件感准张团屋离色" +
		"脸片科倒睛利世刚且由" +
		"送切星导晚表够整认响" +
		"雪流未场该并底深刻平" +
		"伟忙提确近亮轻讲农古" +
		"黑告界拉名呀土清阳照" +
		"办史改历转画造嘴此治" +
		"北必服雨穿内识验传业" +
		"菜爬睡兴形量咱观苦体" +
		"众通冲合破友度术饭公" +
		"旁房极南枪读沙岁线野" +
		"坚空收算至政城劳落钱" +
		"特围弟胜教热展包歌类" +
		"渐强数乡呼性音答哥际" +
		"旧神座章帮啦受系令跳" +
		"非何牛取入岸敢掉忽种" +
		"装顶急林停息句区衣般" +
		"报叶压慢叔背"
)

var (
	// 验证码使用的字体
	DefaultFonts = []string{"chromohv.ttf", "wqy-microhei.ttc", "Flim-Flam.ttf", "ApothecaryFont.ttf", "DENNEthree-dee.ttf", "actionj.ttf", "RitaSmith.ttf"}
)
